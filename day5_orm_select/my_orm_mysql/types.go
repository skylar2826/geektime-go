package my_orm_mysql

import (
	"context"
	"database/sql"
)

type Querier[T any] interface {
	Get(context.Context) (T, error)
	GetMulti(context.Context) ([]*T, error)
}

type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type Query struct {
	SQL  string
	Args []any
}

type QueryBuilder interface {
	Build() (*Query, error)
}
