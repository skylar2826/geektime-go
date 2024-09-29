package my_orm_mysql

import (
	"context"
	"database/sql"
)

type RawQuerier[T any] struct {
	sess Session
	core
	query string
	args  []any
}

func (i *RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  i.query,
		Args: i.args,
	}, nil
}

func RawQuery[T any](sess Session, query string, args ...any) *RawQuerier[T] {
	return &RawQuerier[T]{
		query: query,
		args:  args,
		sess:  sess,
		core:  sess.getCore(),
	}
}

func (i *RawQuerier[T]) ExecContext(ctx context.Context) Result {
	res := Exec(ctx, i.sess, i.core, &QueryContext{
		Type:    "Raw",
		Builder: i,
		Model:   i.model,
	})

	if res.Result != nil {
		return Result{
			res: res.Result.(sql.Result),
		}
	}

	return Result{
		err: res.Err,
	}
}

func (i *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	var err error
	i.model, err = i.R.ParseModel(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, i.sess, i.core, &QueryContext{
		Type:    "Raw",
		Builder: i,
		Model:   i.model,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (i *RawQuerier[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
