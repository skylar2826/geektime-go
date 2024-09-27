package accesslog

import (
	"context"
	"geektime-go/day5_orm/my_orm_mysql"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(query string, args []any)
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(query string, args []any) {
			// args 中如果包含敏感数据的问题没处理
			log.Printf("sql: %s, args: %v", query, args)
		},
	}
}

func (m *MiddlewareBuilder) LogFunc(fn func(query string, args []any)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m *MiddlewareBuilder) Build() my_orm_mysql.Middleware {
	return func(next my_orm_mysql.Handler) my_orm_mysql.Handler {
		return func(ctx context.Context, qc *my_orm_mysql.QueryContext) *my_orm_mysql.QueryResult {
			q, err := qc.Builder.Build()
			if err != nil {
				return &my_orm_mysql.QueryResult{
					Err: err,
				}
			}
			m.logFunc(q.SQL, q.Args)
			res := next(ctx, qc)
			return res
		}
	}
}
