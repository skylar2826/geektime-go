package reflect

import "database/sql"

// DB 我们的db是sql sb的装饰器
type DB struct {
	R  *Register
	DB *sql.DB
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
		R:  NewRegister(),
		DB: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

//func NewDBWithOpts( opts ...DBOpts) (*DB, error) {
//	db := NewDB()
//	for _, opt := range opts {
//		opt(db)
//	}
//	return db, nil
//}
