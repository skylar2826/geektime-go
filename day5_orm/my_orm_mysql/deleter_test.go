package my_orm_mysql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleter(t *testing.T) {
	db, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory", DBWithDialect(DialectMySql))
	require.NoError(t, err)
	testCases := []struct {
		name      string
		entity    any
		wantQuery *Query
		wantErr   error
		builder   QueryBuilder
	}{
		{
			name:    "no from",
			builder: NewDeleter[TestModel](db),
			wantQuery: &Query{
				SQL:  "delete from `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: NewDeleter[TestModel](db).From("`table`"),
			wantQuery: &Query{
				SQL:  "delete from `table`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: NewDeleter[TestModel](db).From(""),
			wantQuery: &Query{
				SQL:  "delete from `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "with db",
			builder: NewDeleter[TestModel](db).From("`test1`.`user`"),
			wantQuery: &Query{
				SQL:  "delete from `test1`.`user`;",
				Args: nil,
			},
		},
		{
			name:    "where",
			builder: NewDeleter[TestModel](db).Where(C("Age").Eq(22)),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where `age` = ?;",
				Args: []any{22},
			},
		},
		{
			name:    "not",
			builder: NewDeleter[TestModel](db).Where(Not(C("Age").Eq(22))),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where  not (`age` = ?);",
				Args: []any{22},
			},
		},
		{
			name:    "and",
			builder: NewDeleter[TestModel](db).Where(C("Age").Eq(22).And(C("Name").Eq("lily"))),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where (`age` = ?) and (`name` = ?);",
				Args: []any{22, "lily"},
			},
		},
		{
			name:    "more and",
			builder: NewDeleter[TestModel](db).Where(C("Age").Eq(22).And(C("Name").Eq("lily").And(C("Sex").Eq(0)))),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where (`age` = ?) and ((`name` = ?) and (`sex` = ?));",
				Args: []any{22, "lily", 0},
			},
		},
		{
			name:    "or",
			builder: NewDeleter[TestModel](db).Where(C("Age").Eq(22).Or(C("Name").Eq("lily"))),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where (`age` = ?) or (`name` = ?);",
				Args: []any{22, "lily"},
			},
		},
		{
			name:    "empty where",
			builder: NewDeleter[TestModel](db).Where(),
			wantQuery: &Query{
				SQL:  "delete from `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "invalid column",
			builder: NewDeleter[TestModel](db).Where(Not(C("XXX").Eq(22))),
			wantErr: fmt.Errorf("field XXX not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, b)
		})
	}
}
