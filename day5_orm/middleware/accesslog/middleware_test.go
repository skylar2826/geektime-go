package accesslog

import (
	"context"
	"database/sql"
	"geektime-go/day5_orm/my_orm_mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestModel struct {
	Id        int
	FirstName string
	Age       int8
	LastName  *sql.NullString
	Name      string
	Sex       int
}

type TestUser struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestMiddleware(t *testing.T) {
	var query string
	var args []any
	m := (&MiddlewareBuilder{}).LogFunc(func(q string, as []any) {
		query = q
		args = as
	})

	db, err := my_orm_mysql.Open("sqlite3", "file:test.db?cache=shared&mode=memory", my_orm_mysql.WithMiddleware(m.Build()))
	require.NoError(t, err)
	_, _ = my_orm_mysql.NewSelector[TestModel](db).Where(my_orm_mysql.C("Id").Eq(1)).Get(context.Background())
	assert.Equal(t, "select * from `test_model` where `id` = ?;", query)
	assert.Equal(t, []any{1}, args)

	_ = my_orm_mysql.NewInsert[TestUser](db).Values(&TestUser{
		Id:        1,
		FirstName: "Tom",
		Age:       18,
		LastName:  &sql.NullString{Valid: true, String: "xi"},
	}).Exec(context.Background())

	assert.Equal(t, "INSERT INTO `test_user`(`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?);", query)
	assert.Equal(t, []any{int64(1), "Tom", int8(18), &sql.NullString{Valid: true, String: "xi"}}, args)
}
