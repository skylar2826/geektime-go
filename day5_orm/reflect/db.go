package reflect

import (
	"database/sql"
	"geektime-go/day5_orm/internal/valuer"
	"geektime-go/day5_orm/model"
)

// DB 我们的db是sql sb的装饰器
type DB struct {
	R       *model.Register
	DB      *sql.DB
	Creator valuer.Creator
}

type DBOpts func(db *DB)

func Open(driver string, datasourceName string, opts ...DBOpts) (*DB, error) {
	db, err := sql.Open(driver, datasourceName)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

// OpenDB 支持单独传入db, 也方便单元测试
func OpenDB(db *sql.DB, opts ...DBOpts) (*DB, error) {
	res := &DB{
		R:       model.NewRegister(),
		DB:      db,
		Creator: valuer.NewUnsafeValue,
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
