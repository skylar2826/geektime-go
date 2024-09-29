package my_orm_mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"geektime-go/day5_orm/internal"
	"github.com/DATA-DOG/go-sqlmock"
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

func TestSelector(t *testing.T) {
	db, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory", DBWithDialect(DialectMySql))
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
		//{
		//	name:    "from",
		//	builder: NewSelector[TestModel](db).From("`table`"),
		//	wantQuery: &Query{
		//		SQL:  "select * from `table`;",
		//		Args: nil,
		//	},
		//},
		//{
		//	name:    "empty from",
		//	builder: NewSelector[TestModel](db).From(""),
		//	wantQuery: &Query{
		//		SQL:  "select * from `test_model`;",
		//		Args: nil,
		//	},
		//},
		//{
		//	name:    "with db",
		//	builder: NewSelector[TestModel](db).From("`test1`.`user`"),
		//	wantQuery: &Query{
		//		SQL:  "select * from `test1`.`user`;",
		//		Args: nil,
		//	},
		//},
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
	db, err := OpenDB(mockDB)
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
	db, err := OpenDB(mockDB)
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
	db, err := OpenDB(mockDB)
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
		{
			// 必须是前面存在的列
			name: "group by",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("Age")).GroupBy(C("FirstName")),
			wantRes: &Query{
				SQL: "select `first_name`,`age` from `test_model` group by `first_name`;",
			},
		},

		{
			name: "group by more",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("Age"), C("LastName")).GroupBy(C("FirstName"), C("Age")),
			wantRes: &Query{
				SQL: "select `first_name`,`age`,`last_name` from `test_model` group by `first_name`,`age`;",
			},
		},
		// having要和group by一起出现，不能单独出现；可以考虑单独使用having报错
		{
			name: "having predicate",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("Age"), C("LastName")).GroupBy(C("FirstName"), C("Age")).Having(Avg("age").Lt(18)),
			wantRes: &Query{
				SQL:  "select `first_name`,`age`,`last_name` from `test_model` group by `first_name`,`age` having AVG(`age`) < ?;",
				Args: []any{18},
			},
		},
		{
			name: "having predicate more",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("Age"), C("LastName")).GroupBy(C("FirstName"), C("Age")).Having(C("Age").Lt(18).And(C("FirstName").Eq("lili"))),
			wantRes: &Query{
				SQL:  "select `first_name`,`age`,`last_name` from `test_model` group by `first_name`,`age` having (`age` < ?) and (`first_name` = ?);",
				Args: []any{18, "lili"},
			},
		},
		{
			name: "order by",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("Age")).orderBy(ASC(C("Age"))),
			wantRes: &Query{
				SQL: "select `first_name`,`age` from `test_model` order by `age` ASC;",
			},
		},
		{
			name: "order by more",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("Age")).orderBy(ASC(C("Age")), DESC(C(`FirstName`))),
			wantRes: &Query{
				SQL: "select `first_name`,`age` from `test_model` order by `age` ASC,`first_name` DESC;",
			},
		},
		{
			name: "limit",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("Age")).orderBy(ASC(C("Age")), DESC(C(`FirstName`))).Limit(1),
			wantRes: &Query{
				SQL:  "select `first_name`,`age` from `test_model` order by `age` ASC,`first_name` DESC limit ?;",
				Args: []any{1},
			},
		},
		{
			name: "offset",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("Age")).orderBy(ASC(C("Age")), DESC(C(`FirstName`))).Limit(1).Offset(2),
			wantRes: &Query{
				SQL:  "select `first_name`,`age` from `test_model` order by `age` ASC,`first_name` DESC limit ? offset ?;",
				Args: []any{1, 2},
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

func TestSelect_Join(t *testing.T) {
	db := MemoryDB(t)

	type Order struct {
		Id        int
		UsingCol1 string
		UsingCol2 string
	}

	type OrderDetail struct {
		OrderId   int
		ItemId    int
		UsingCol1 string
		UsingCol2 string
	}

	type Item struct {
		Id int
	}

	testCases := []struct {
		name      string
		s         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "specify table",
			s:    NewSelector[Order](db).From(TableOf(&OrderDetail{})),
			wantQuery: &Query{
				SQL: "select * from `order_detail`;",
			},
		},

		{
			name: "join using",
			s: func() QueryBuilder {
				t1 := TableOf(&OrderDetail{})
				t2 := TableOf(&Order{})
				t3 := t1.join(t2).Using("UsingCol1", "UsingCol2")
				return NewSelector[Order](db).From(t3)
			}(),
			wantQuery: &Query{
				SQL: "select * from (`order_detail` Join `order` Using (`using_col1`,`using_col2`));",
			},
		},
		{
			name: "join on",
			s: func() QueryBuilder {
				t1 := TableOf(&OrderDetail{})
				t2 := TableOf(&Order{})
				// NewSelector[Order](db)中的model是Order; 而复杂查询中联合了多个元数据Order、OrderDetail，并操作Order.Id、OrderDetail.OrderId
				// 所以需要改在buildColumn, 使用每个entity自己的元数据
				t3 := t1.join(t2).On(t2.C("Id").Eq(t1.C("OrderId")))
				return NewSelector[Order](db).From(t3)
			}(),
			wantQuery: &Query{
				SQL: "select * from (`order_detail` Join `order` On (`id` = `order_id`));",
			},
		},
		{
			name: "join on alias",
			s: func() QueryBuilder {
				t1 := TableOf(&OrderDetail{}).As("t1")
				t2 := TableOf(&Order{}).As("t2")
				// NewSelector[Order](db)中的model是Order; 而复杂查询中联合了多个元数据Order、OrderDetail，并操作Order.Id、OrderDetail.OrderId
				// 所以需要改在buildColumn, 使用每个entity自己的元数据
				t3 := t1.join(t2).On(t2.C("Id").Eq(t1.C("OrderId")))
				return NewSelector[Order](db).From(t3)
			}(),
			wantQuery: &Query{
				SQL: "select * from (`order_detail` As `t1` Join `order` As `t2` On (`t2`.`id` = `t1`.`order_id`));",
			},
		},
		{
			name: "join join on",
			s: func() QueryBuilder {
				t1 := TableOf(&OrderDetail{}).As("t1")
				t2 := TableOf(&Order{}).As("t2")
				t3 := t1.join(t2).On(t2.C("Id").Eq(t1.C("OrderId")))
				t4 := TableOf(&Item{}).As("t4")
				t5 := t3.join(t4).On(t1.C("ItemId").Eq(t4.C("Id")))
				return NewSelector[Order](db).From(t5)
			}(),
			wantQuery: &Query{
				SQL: "select * from ((`order_detail` As `t1` Join `order` As `t2` On (`t2`.`id` = `t1`.`order_id`)) Join `item` As `t4` On (`t1`.`item_id` = `t4`.`id`));",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.s.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}
