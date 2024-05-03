package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) DB {
	return DB{db: db}
}

func (db DB) Tx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
