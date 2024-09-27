package my_orm_mysql

import (
	"context"
	"database/sql"
	"fmt"
	"geektime-go/day5_orm/internal/valuer"
	"geektime-go/day5_orm/model"
)

var (
	_ Session = &DB{}
)

// DB 我们的db是sql sb的装饰器
type DB struct {
	DB *sql.DB
	core
}

type DBOpts func(db *DB)

func (db *DB) getCore() core {
	return db.core
}

func WithMiddleware(mdls ...Middleware) DBOpts {
	return func(db *DB) {
		db.middlewares = mdls
	}
}

func Open(driver string, datasourceName string, opts ...DBOpts) (*DB, error) {
	db, err := sql.Open(driver, datasourceName)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

func DBWithDialect(dialect Dialect) DBOpts {
	return func(db *DB) {
		db.dialect = dialect
	}
}

// OpenDB 支持单独传入db, 也方便单元测试
func OpenDB(db *sql.DB, opts ...DBOpts) (*DB, error) {
	res := &DB{
		core: core{
			R:       model.NewRegister(),
			Creator: valuer.NewUnsafeValue,
			dialect: DialectMySql,
		},

		DB: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBUseReflect() DBOpts {
	return func(db *DB) {
		db.Creator = valuer.NewReflectValue
	}
}

//func NewDBWithOpts( opts ...DBOpts) (*DB, error) {
//	db := NewDB()
//	for _, opt := range opts {
//		opt(db)
//	}
//	return db, nil
//}

func (db *DB) BeginTx(ctx context.Context, txOpts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx: tx,
		db: db,
	}, nil
}

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.DB.ExecContext(ctx, query, args...)
}

func (db *DB) DoTx(ctx context.Context, fn func(ctx context.Context, tx *Tx) error, txOpts *sql.TxOptions) (err error) {
	var tx *Tx
	tx, err = db.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}

	panicked := true
	err = fn(ctx, tx)
	panicked = false

	defer func() {
		if err != nil || panicked {
			e := tx.Rollback()

			err = fmt.Errorf("业务错误: %w, 回滚错误：%s, 是否panic: %t", err, e, panicked)
		} else {
			err = tx.Commit()
		}
	}()
	return err
}
