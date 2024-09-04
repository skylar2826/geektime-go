package day5_orm

import (
	"geektime-go/day5_orm/orm_mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func RunORMSelect() {
	//base_mysql.RunBase()
	orm_mysql.RunOrm()
}
