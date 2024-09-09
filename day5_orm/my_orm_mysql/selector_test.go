package my_orm_mysql

import (
	"database/sql"
	"fmt"
	rft "geektime-go/day5_orm/reflect"
	"github.com/stretchr/testify/assert"
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

func TestSelector(t *testing.T) {
	db := rft.NewDB()
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no from",
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "select * from `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: NewSelector[TestModel](db).From("`table`"),
			wantQuery: &Query{
				SQL:  "select * from `table`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL:  "select * from `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "with db",
			builder: NewSelector[TestModel](db).From("`test1`.`user`"),
			wantQuery: &Query{
				SQL:  "select * from `test1`.`user`;",
				Args: nil,
			},
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](db).Where(C("Age").Eq(22)),
			wantQuery: &Query{
				SQL:  "select * from `test_model` where `age` = ?;",
				Args: []any{22},
			},
		},
		{
			name:    "not",
			builder: NewSelector[TestModel](db).Where(Not(C("Age").Eq(22))),
			wantQuery: &Query{
				SQL:  "select * from `test_model` where  not (`age` = ?);",
				Args: []any{22},
			},
		},
		{
			name:    "and",
			builder: NewSelector[TestModel](db).Where(C("Age").Eq(22).And(C("Name").Eq("lily"))),
			wantQuery: &Query{
				SQL:  "select * from `test_model` where (`age` = ?) and (`name` = ?);",
				Args: []any{22, "lily"},
			},
		},
		{
			name:    "more and",
			builder: NewSelector[TestModel](db).Where(C("Age").Eq(22).And(C("Name").Eq("lily").And(C("Sex").Eq(0)))),
			wantQuery: &Query{
				SQL:  "select * from `test_model` where (`age` = ?) and ((`name` = ?) and (`sex` = ?));",
				Args: []any{22, "lily", 0},
			},
		},
		{
			name:    "or",
			builder: NewSelector[TestModel](db).Where(C("Age").Eq(22).Or(C("Name").Eq("lily"))),
			wantQuery: &Query{
				SQL:  "select * from `test_model` where (`age` = ?) or (`name` = ?);",
				Args: []any{22, "lily"},
			},
		},
		{
			name:    "empty where",
			builder: NewSelector[TestModel](db).Where(),
			wantQuery: &Query{
				SQL:  "select * from `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "invalid column",
			builder: NewSelector[TestModel](db).Where(Not(C("XXX").Eq(22))),
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
