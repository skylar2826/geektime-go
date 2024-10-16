package testdata

import (
	"database/sql"
	sqlx "database/sql"
)

type User struct {
	Name     string
	Age      int
	NickName *sql.NullString
	Phone    *sqlx.NullString
	Picture  []byte
}

type UserDetail struct {
	Address string
}
