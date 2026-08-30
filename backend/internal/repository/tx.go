package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBTX is a common interface implemented by *pgxpool.Pool and pgx.Tx
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
}

type txKey struct{}

// TransactionManager allows executing operations inside a transaction
type TransactionManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type PgTransactionManager struct {
	pool *pgxpool.Pool
}

func NewPgTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return &PgTransactionManager{pool: pool}
}

// RunInTx starts a transaction, stores it in the context, runs the function, and commits/rolls back.
func (m *PgTransactionManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// If a transaction is already in context, just use it (nested transactions not fully supported without savepoints, but allows reuse)
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetDB extracts the database connection (either a Tx or a Pool) from the context.
func GetDB(ctx context.Context, pool *pgxpool.Pool) DBTX {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}
