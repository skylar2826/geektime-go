//go:build e2e

package integration

import (
	"context"
	"fmt"
	"geektime-go/day5_orm/my_orm_mysql"
	day5_orm_select "geektime-go/day5_orm/types"
	"geektime-go/day7/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type suites struct {
	suite.Suite
	driver string
	dsn    string
	db     *my_orm_mysql.DB
}

func (i *suites) SetupSuite() {
	db, err := my_orm_mysql.Open(i.driver, i.dsn)
	require.NoError(i.T(), err)
	err = db.Wait() // 等待db连接好了再跑集成测试
	require.NoError(i.T(), err)
	i.db = db
}

type InsertSuite struct {
	suites
}

func TestMySqlInsert(t *testing.T) {
	datasourceName := fmt.Sprint(day5_orm_select.UserName, ":", day5_orm_select.Password, "@tcp(", day5_orm_select.Ip, ":", day5_orm_select.Port, ")/", day5_orm_select.DbName)
	suite.Run(t, &InsertSuite{
		suites: suites{
			driver: "mysql",
			dsn:    datasourceName,
		},
	})
}

func (i *InsertSuite) TestInsert() {
	db := i.db
	t := i.T()

	testCases := []struct {
		name string
		i    *my_orm_mysql.Insert[test.SimpleStruct]
		//wantRes *my_orm_mysql.QueryResult
		wantAffectId int64 // 插入行数
	}{
		{
			name:         "insert one",
			i:            my_orm_mysql.NewInsert[test.SimpleStruct](db).Values(test.NewSimpleStruct(12)),
			wantAffectId: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			res := tc.i.Exec(ctx)
			affectId, err := res.RowsAffected()
			assert.NoError(t, err)
			assert.Equal(t, affectId, tc.wantAffectId)
		})
	}
}

//
//func TestMysqlInsert(t *testing.T) {
//	datasourceName := fmt.Sprint(day5_orm_select.UserName, ":", day5_orm_select.Password, "@tcp(", day5_orm_select.Ip, ":", day5_orm_select.Port, ")/", day5_orm_select.DbName)
//	testInsert(t, "mysql", datasourceName)
//}
//
//func testInsert(t *testing.T, driver string, dsn string) {
//	db, err := my_orm_mysql.Open("mysql", dsn)
//	require.NoError(t, err)
//	testCases := []struct {
//		name string
//		i    *my_orm_mysql.Insert[test.SimpleStruct]
//		//wantRes *my_orm_mysql.QueryResult
//		wantAffectId int64 // 插入行数
//	}{
//		{
//			name:         "insert one",
//			i:            my_orm_mysql.NewInsert[test.SimpleStruct](db).Values(test.NewSimpleStruct(12)),
//			wantAffectId: 1,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//			defer cancel()
//			res := tc.i.Exec(ctx)
//			var affectId int64
//			affectId, err = res.RowsAffected()
//			assert.NoError(t, err)
//			assert.Equal(t, affectId, tc.wantAffectId)
//		})
//	}
//}
