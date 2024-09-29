package my_orm_mysql

import (
	"context"
	"geektime-go/day5_orm/model"
)

type QueryContext struct {
	// 查询类型
	Type    string
	Builder QueryBuilder
	Model   *model.Model
}

type QueryResult struct {
	// selector 中是 *T []*T
	// insert 中是 Result
	Result any
	Err    error
}

type Handler func(ctx context.Context, queryCtx *QueryContext) *QueryResult

type Middleware func(next Handler) Handler
