package my_orm_mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"geektime-go/day5_orm/internal"
	rft "geektime-go/day5_orm/reflect"
	"github.com/DATA-DOG/go-sqlmock"
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

func TestSelector(t *testing.T) {
	db, err := rft.Open("", "")
	require.NoError(t, err)
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

// 要全部一起跑，不然sqlmock配对有问题
func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := rft.OpenDB(mockDB)
	require.NoError(t, err)

	// mock与测试一对一，顺序不可变

	// invalid query
	mock.ExpectQuery("select .*").WillReturnError(errors.New("query error"))

	// now rows
	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	mock.ExpectQuery("select .* from `test_model` where .*").WillReturnRows(rows)

	// row 要新建一个，直接用上面的rows会有影响
	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "Tom", "18", "lili")
	mock.ExpectQuery("select .* from `test_model` where .*").WillReturnRows(rows)

	testCases := []struct {
		name    string
		s       *Selector[TestModel]
		wantErr error
		wantRes *TestModel
	}{
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(C("xxx").Eq(1)),
			wantErr: errors.New("field xxx not found"),
		},
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(9999)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "now rows",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(2)),
			wantErr: internal.ErrorNoRows,
		},
		{
			name: "row",
			s:    NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantRes: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "lili"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.GetV1(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
