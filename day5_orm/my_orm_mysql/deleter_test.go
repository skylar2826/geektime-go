package my_orm_mysql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleter(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantQuery *Query
		wantErr   error
		builder   QueryBuilder
	}{
		{
			name:    "no from",
			builder: &Deleter[TestModel]{},
			wantQuery: &Query{
				SQL:  "delete from `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: (&Deleter[TestModel]{}).From("`table`"),
			wantQuery: &Query{
				SQL:  "delete from `table`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: (&Deleter[TestModel]{}).From(""),
			wantQuery: &Query{
				SQL:  "delete from `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "with db",
			builder: (&Deleter[TestModel]{}).From("`test1`.`user`"),
			wantQuery: &Query{
				SQL:  "delete from `test1`.`user`;",
				Args: nil,
			},
		},
		{
			name:    "where",
			builder: (&Deleter[TestModel]{}).Where(C("Age").Eq(22)),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where `age` = ?;",
				Args: []any{22},
			},
		},
		{
			name:    "not",
			builder: (&Deleter[TestModel]{}).Where(Not(C("Age").Eq(22))),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where  not (`age` = ?);",
				Args: []any{22},
			},
		},
		{
			name:    "and",
			builder: (&Deleter[TestModel]{}).Where(C("Age").Eq(22).And(C("Name").Eq("lily"))),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where (`age` = ?) and (`name` = ?);",
				Args: []any{22, "lily"},
			},
		},
		{
			name:    "more and",
			builder: (&Deleter[TestModel]{}).Where(C("Age").Eq(22).And(C("Name").Eq("lily").And(C("Sex").Eq(0)))),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where (`age` = ?) and ((`name` = ?) and (`sex` = ?));",
				Args: []any{22, "lily", 0},
			},
		},
		{
			name:    "or",
			builder: (&Deleter[TestModel]{}).Where(C("Age").Eq(22).Or(C("Name").Eq("lily"))),
			wantQuery: &Query{
				SQL:  "delete from `test_model` where (`age` = ?) or (`name` = ?);",
				Args: []any{22, "lily"},
			},
		},
		{
			name:    "empty where",
			builder: (&Deleter[TestModel]{}).Where(),
			wantQuery: &Query{
				SQL:  "delete from `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "invalid column",
			builder: (&Deleter[TestModel]{}).Where(Not(C("XXX").Eq(22))),
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
