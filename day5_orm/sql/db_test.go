package sql

import (
	"context"
	"database/sql"
	"geektime-go/day5_orm/types"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	// shared 意味着有很多goroutine可以同时操作db
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	err = db.Ping()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)`)
	require.NoError(t, err)

	var res sql.Result
	res, err = db.ExecContext(ctx, "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES (?, ?, ?, ?)", 1, "Tom", 18, "xi")
	require.NoError(t, err)

	var rowIdx int64
	rowIdx, err = res.RowsAffected()
	require.NoError(t, err)

	var id int64
	id, err = res.LastInsertId()
	require.NoError(t, err)
	log.Println("affected:", rowIdx, " insert id:", id)

	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, `select * from test_model where id=?`, 1)

	resRows := make([]interface{}, 0, 4)
	// QueryContext 可能有多条数据，可能没有数据
	if rows != nil {
		for rows.Next() {
			resRow := &types.TestModel{}
			err = rows.Scan(&resRow.Id, &resRow.FirstName, &resRow.Age, &resRow.LastName)
			resRows = append(resRows, resRow)
		}
	}

	var row *sql.Row
	row = db.QueryRowContext(ctx, `select * from test_model where id=?`, 1)
	require.NoError(t, row.Err())
	resRow := &types.TestModel{}
	err = row.Scan(&resRow.Id, &resRow.FirstName, &resRow.Age, &resRow.LastName)
	require.NoError(t, err)

	// QueryRowContext 预期有一条数据，没有数据则报错
	row = db.QueryRowContext(ctx, `select * from test_model where id=?`, 2)
	require.NoError(t, row.Err())
	resRow = &types.TestModel{}
	err = row.Scan(&resRow.Id, &resRow.FirstName, &resRow.LastName)
	require.Error(t, err, sql.ErrNoRows)
	cancel()

}
