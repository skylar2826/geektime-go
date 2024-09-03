package day5_orm_select

import (
	"geektime-go/day5_orm_select/orm_mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func RunORMSelect() {
	//base_mysql.RunBase()
	orm_mysql.RunOrm()
}
