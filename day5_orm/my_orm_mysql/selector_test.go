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
	db, err := rft.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
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
			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

// 要全部一起跑，不然sqlmock配对有问题
func TestSelector_GetMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := rft.OpenDB(mockDB)
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "Tom", "18", "lili")
	rows.AddRow("2", "xx", "8", "sss")
	mock.ExpectQuery("select .* from `test_model` where .*").WillReturnRows(rows)

	testCases := []struct {
		name    string
		s       *Selector[TestModel]
		wantErr error
		wantRes []*TestModel
	}{

		{
			name: "muti row",
			s:    NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantRes: []*TestModel{&TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "lili"},
			}, &TestModel{
				Id:        2,
				FirstName: "xx",
				Age:       8,
				LastName:  &sql.NullString{Valid: true, String: "sss"},
			}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.GetMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestSelector_Select(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	db, err := rft.OpenDB(mockDB)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		s       *Selector[TestModel]
		wantRes *Query
		wantErr error
	}{
		{
			name:    "invalid",
			s:       NewSelector[TestModel](db).Select(C("invalid")),
			wantErr: fmt.Errorf("field invalid not found"),
		},
		{
			name: "select",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("Age")),
			wantRes: &Query{
				SQL: "select `first_name`,`age` from `test_model`;",
			},
		},
		{
			name: "select alias",
			s:    NewSelector[TestModel](db).Select(C("FirstName").AS("AAA"), C("Age")),
			wantRes: &Query{
				SQL: "select `first_name` AS `AAA`,`age` from `test_model`;",
			},
		},
		{
			name: "select *",
			s:    NewSelector[TestModel](db),
			wantRes: &Query{
				SQL: "select * from `test_model`;",
			},
		},
		{
			name: "select avg",
			s:    NewSelector[TestModel](db).Select(Avg("FirstName")),
			wantRes: &Query{
				SQL: "select AVG(`FirstName`) from `test_model`;",
			},
		},
		{
			name: "select avg",
			s:    NewSelector[TestModel](db).Select(Avg("FirstName").AS("AAA")),
			wantRes: &Query{
				SQL: "select AVG(`FirstName`) AS `AAA` from `test_model`;",
			},
		},
		{
			name: "select multi avg",
			s:    NewSelector[TestModel](db).Select(Avg("FirstName"), Avg("age")),
			wantRes: &Query{
				SQL: "select AVG(`FirstName`),AVG(`age`) from `test_model`;",
			},
		},
		{
			name: "select rawExpr",
			s:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT `first_name`)")),
			wantRes: &Query{
				SQL: "select COUNT(DISTINCT `first_name`) from `test_model`;",
			},
		},
		{
			name: "where rawExpr",
			s:    NewSelector[TestModel](db).Where(Raw("id < ?", 18).AsPredicate()),
			wantRes: &Query{
				SQL:  "select * from `test_model` where (id < ?);",
				Args: []any{18},
			},
		},
		{
			name: "where rawExpr used in predicate",
			s:    NewSelector[TestModel](db).Where(C("Id").Eq(Raw("`age` < ?", 18))),
			wantRes: &Query{
				SQL:  "select * from `test_model` where `id` = (`age` < ?);",
				Args: []any{18},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run(tc.name, func(t *testing.T) {
				res, err := tc.s.Build()
				assert.Equal(t, tc.wantErr, err)
				if err != nil {
					return
				}
				assert.Equal(t, tc.wantRes, res)
			})
		})
	}
}
