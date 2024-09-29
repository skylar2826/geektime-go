package my_orm_mysql

import (
	"context"
	"database/sql"
)

// Session 会话/上下文/分组
type Session interface {
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error) // sql.Result
}
