package config

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB creates a pgxpool connection pool and verifies connectivity.
func ConnectDB(cfg *Config) (*pgxpool.Pool, error) {
	dsn := cfg.DatabaseURL()

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	slog.Info("database connected", "host", cfg.DBHost, "db", cfg.DBName)
	return pool, nil
}

// RunMigrations reads SQL files from migrationsFS (an embedded filesystem,
// see backend/migrations/embed.go) and executes each ".up.sql" file that has
// not yet been applied, in filename order. Each migration runs inside its
// own transaction: a failure rolls the whole file back and is never marked
// as applied, so partial schema changes can't get silently "recorded" as
// done. An advisory lock also prevents two instances from racing to apply
// migrations concurrently on startup.
func RunMigrations(pool *pgxpool.Pool, migrationsFS fs.FS) error {
	ctx := context.Background()

	// Create tracking table if not exists.
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	// Advisory lock: avoids two processes/replicas applying migrations at
	// the same time and stepping on each other.
	const migrationLockID = 727384910
	if _, err := pool.Exec(ctx, "SELECT pg_advisory_lock($1)", migrationLockID); err != nil {
		return fmt.Errorf("acquire migration lock: %w", err)
	}
	defer func() {
		if _, err := pool.Exec(ctx, "SELECT pg_advisory_unlock($1)", migrationLockID); err != nil {
			slog.Warn("failed to release migration advisory lock", "error", err)
		}
	}()

	// Find all .up.sql files in the embedded FS.
	entries, err := fs.ReadDir(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	if len(upFiles) == 0 {
		slog.Warn("no migration files found in embedded FS — check the go:embed directive picked them up")
	}

	for _, fileName := range upFiles {
		version := strings.TrimSuffix(fileName, ".up.sql")

		// Check if already applied.
		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", version, err)
		}
		if exists {
			slog.Debug("migration already applied", "version", version)
			continue
		}

		sqlBytes, err := fs.ReadFile(migrationsFS, fileName)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", fileName, err)
		}

		if err := applyMigration(ctx, pool, version, string(sqlBytes)); err != nil {
			return fmt.Errorf("execute migration %s: %w", fileName, err)
		}

		slog.Info("migration applied", "version", version)
	}

	return nil
}

// applyMigration executes a single migration file's SQL and records it as
// applied, both inside one transaction. If anything fails, the transaction
// rolls back automatically (via the deferred Rollback, which is a no-op
// after a successful Commit) — so a migration is either fully applied and
// recorded, or neither happens.
func applyMigration(ctx context.Context, pool *pgxpool.Pool, version string, sqlText string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, sqlText); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return fmt.Errorf("record migration: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
