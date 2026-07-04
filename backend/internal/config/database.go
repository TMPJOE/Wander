package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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

// RunMigrations reads SQL files from the migrations directory and executes them.
// It uses a simple `schema_migrations` table to track applied migrations.
func RunMigrations(pool *pgxpool.Pool, migrationsDir string) error {
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

	// Find all .up.sql files.
	entries, err := os.ReadDir(migrationsDir)
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

		// Read and execute.
		sqlBytes, err := os.ReadFile(filepath.Join(migrationsDir, fileName))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", fileName, err)
		}

		_, err = pool.Exec(ctx, string(sqlBytes))
		if err != nil {
			// If the migration fails because the object already exists (e.g. column added manually),
			// consider the migration effectively applied and record it so subsequent runs skip it.
			if strings.Contains(err.Error(), "already exists") {
				slog.Warn("migration execution reported existing object; recording as applied", "version", version, "error", err)
				_, insErr := pool.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version)
				if insErr != nil {
					return fmt.Errorf("record migration %s after partial apply: %w", version, insErr)
				}
				continue
			}
			return fmt.Errorf("execute migration %s: %w", fileName, err)
		}

		// Record as applied.
		_, err = pool.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version)
		if err != nil {
			return fmt.Errorf("record migration %s: %w", version, err)
		}

		slog.Info("migration applied", "version", version)
	}

	return nil
}
