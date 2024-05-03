package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Querier interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func ChooseQuerier(tx *sqlx.Tx, db *sqlx.DB) Querier {
	if tx != nil {
		return tx
	}

	return db
}
