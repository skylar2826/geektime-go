package reflect

type DB struct {
	R *Register
}

func NewDB() *DB {
	return &DB{
		R: NewRegister(),
	}
}

type DBOpts func(db *DB)

func NewDBWithOpts(opts ...DBOpts) (*DB, error) {
	db := NewDB()
	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}
