//go:build e2e

package integration

import (
	"context"
	"fmt"
	"geektime-go/day5_orm/my_orm_mysql"
	day5_orm_select "geektime-go/day5_orm/types"
	"geektime-go/day7/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SelectSuite struct {
	suites
}

func (s *SelectSuite) SetupSuite() {
	fmt.Println("1.....")
	s.suites.SetupSuite()
	res := my_orm_mysql.NewInsert[test.SimpleStruct](s.db).Values(test.NewSimpleStruct(13)).Exec(context.Background())
	id, err := res.LastInsertId()
	fmt.Println(id, err)
	var rowId int64
	rowId, err = res.RowsAffected()
	fmt.Println(rowId, err)
	fmt.Println("-----------------------------------------")
}

func TestMySqlSelect(t *testing.T) {
	fmt.Println("2.....")
	datasourceName := fmt.Sprint(day5_orm_select.UserName, ":", day5_orm_select.Password, "@tcp(", day5_orm_select.Ip, ":", day5_orm_select.Port, ")/", day5_orm_select.DbName)
	suite.Run(t, &SelectSuite{
		suites: suites{
			driver: "mysql",
			dsn:    datasourceName,
		},
	})
}

func (i *SelectSuite) TestSelect() {
	fmt.Println("3.....")
	db := i.db
	t := i.T()

	testCases := []struct {
		name string
		i    *my_orm_mysql.Selector[test.SimpleStruct]
		//wantRes *my_orm_mysql.QueryResult
		wantRes *test.SimpleStruct
	}{
		{
			name:    "insert one",
			i:       my_orm_mysql.NewSelector[test.SimpleStruct](db).Where(my_orm_mysql.C("Id").Eq(13)),
			wantRes: test.NewSimpleStruct(13),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			res, err := tc.i.Get(ctx)
			assert.NoError(t, err)
			assert.Equal(t, res, tc.wantRes)
		})
	}
}
