package my_orm_mysql

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Id        int
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestSelector(t *testing.T) {
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no from",
			builder: &selector[TestModel]{},
			wantQuery: &Query{
				SQL:  "select * from `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: (&selector[TestModel]{}).From("`table`"),
			wantQuery: &Query{
				SQL:  "select * from `table`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: (&selector[TestModel]{}).From(""),
			wantQuery: &Query{
				SQL:  "select * from `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "with db",
			builder: (&selector[TestModel]{}).From("`test1`.`user`"),
			wantQuery: &Query{
				SQL:  "select * from `test1`.`user`;",
				Args: nil,
			},
		},
		{
			name:    "where",
			builder: (&selector[TestModel]{}).Where(C("Age").Eq(22)),
			wantQuery: &Query{
				SQL:  "select * from `TestModel` where `Age` = ?;",
				Args: []any{22},
			},
		},
		{
			name:    "not",
			builder: (&selector[TestModel]{}).Where(Not(C("Age").Eq(22))),
			wantQuery: &Query{
				SQL:  "select * from `TestModel` where  not (`Age` = ?);",
				Args: []any{22},
			},
		},
		{
			name:    "and",
			builder: (&selector[TestModel]{}).Where(C("Age").Eq(22).And(C("Name").Eq("lily"))),
			wantQuery: &Query{
				SQL:  "select * from `TestModel` where (`Age` = ?) and (`Name` = ?);",
				Args: []any{22, "lily"},
			},
		},
		{
			name:    "more and",
			builder: (&selector[TestModel]{}).Where(C("Age").Eq(22).And(C("Name").Eq("lily").And(C("Sex").Eq(0)))),
			wantQuery: &Query{
				SQL:  "select * from `TestModel` where (`Age` = ?) and ((`Name` = ?) and (`Sex` = ?));",
				Args: []any{22, "lily", 0},
			},
		},
		{
			name:    "or",
			builder: (&selector[TestModel]{}).Where(C("Age").Eq(22).Or(C("Name").Eq("lily"))),
			wantQuery: &Query{
				SQL:  "select * from `TestModel` where (`Age` = ?) or (`Name` = ?);",
				Args: []any{22, "lily"},
			},
		},
		{
			name:    "empty where",
			builder: (&selector[TestModel]{}).Where(),
			wantQuery: &Query{
				SQL:  "select * from `TestModel`;",
				Args: nil,
			},
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
