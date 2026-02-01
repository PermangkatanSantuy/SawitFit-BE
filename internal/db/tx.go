package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

func (db *DB) WithUser(
	ctx context.Context,
	userID string,
	fn func(tx pgx.Tx) error,
) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// RLS bridge
	_, err = tx.Exec(
		ctx,
		"SELECT set_config('request.jwt.claim.sub', $1, true)",
		userID,
	)
	if err != nil {
		return fmt.Errorf("set jwt claim: %w", err)
	}
	
	// DB logic
	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}