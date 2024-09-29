package my_orm_mysql

import (
	"context"
	"database/sql"
	"geektime-go/day5_orm/internal"
	"geektime-go/day5_orm/internal/valuer"
	"geektime-go/day5_orm/model"
)

type core struct {
	model       *model.Model
	dialect     Dialect
	Creator     valuer.Creator
	R           *model.Register
	middlewares []Middleware
}

func getHandler[T any](ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}

	var rows *sql.Rows
	rows, err = sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}

	if !rows.Next() {
		return &QueryResult{
			Err: internal.ErrorNoRows,
		}
	}

	tp := new(T)
	v := c.Creator(qc.Model, tp)
	err = v.SetColumns(rows)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	return &QueryResult{
		Result: tp,
	}

}

func get[T any](ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
	var root Handler = func(ctx context.Context, queryCtx *QueryContext) *QueryResult {
		return getHandler[T](ctx, sess, c, queryCtx)
	}
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		root = c.middlewares[i](root)
	}
	res := root(ctx, qc)
	if res.Result != nil {
		return &QueryResult{
			Result: res.Result.(*T),
			Err:    res.Err,
		}
	}
	return &QueryResult{
		Err: res.Err,
	}
}

func getMuti[T any](ctx context.Context, sess Session, c core, qc *QueryContext) ([]*T, error) {
	q, err := qc.Builder.Build()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	var tps []*T
	for rows.Next() {
		tp := new(T)
		v := c.Creator(c.model, tp)
		err = v.SetColumns(rows)
		if err != nil {
			return nil, err
		}

		tps = append(tps, tp)
	}
	return tps, nil
}

func ExecHandler(ctx context.Context, sess Session, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{Err: err}
	}
	var res sql.Result
	res, err = sess.execContext(ctx, q.SQL, q.Args...)
	return &QueryResult{
		Result: res,
		Err:    err,
	}
}

func Exec(ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {

	var root Handler = func(ctx context.Context, queryCtx *QueryContext) *QueryResult {
		return ExecHandler(ctx, sess, queryCtx)
	}
	for j := len(c.middlewares) - 1; j >= 0; j-- {
		root = c.middlewares[j](root)
	}

	return root(ctx, qc)

}
