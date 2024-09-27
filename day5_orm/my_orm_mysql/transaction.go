package my_orm_mysql

import (
	"context"
	"database/sql"
	"errors"
)

var (
	_ Session = &Tx{}
)

type Tx struct {
	tx *sql.Tx
	db *DB
}

func (tx *Tx) getCore() core {
	return tx.db.getCore()
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}
func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) RollbackIfNotCommit() error {
	err := tx.tx.Rollback()
	if errors.Is(err, sql.ErrTxDone) {
		return nil
	}
	return err
}

func (tx *Tx) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

func (tx *Tx) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, args...)
}
