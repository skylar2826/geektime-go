package my_orm_mysql

import (
	"context"
)

type Querier[T any] interface {
	Get(context.Context) (T, error)
	GetMulti(context.Context) ([]*T, error)
}

type Executor interface {
	ExecContext(ctx context.Context) Result
}

type Query struct {
	SQL  string
	Args []any
}

type QueryBuilder interface {
	Build() (*Query, error)
}
