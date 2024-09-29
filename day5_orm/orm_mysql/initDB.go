package orm_mysql

import (
	"fmt"
	day5_orm_select "geektime-go/day5_orm/types"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func initDB() error {
	orm.RegisterModel(new(day5_orm_select.User))
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		return err
	}
	//root:wxFqogsjy5+V@tcp(127.0.0.1:3306)/test
	//root:15271908767Aa!@tcp(127.0.0.1:3306)/test
	datasourceName := fmt.Sprint(day5_orm_select.UserName, ":", day5_orm_select.Password, "@tcp(", day5_orm_select.Ip, ":", day5_orm_select.Port, ")/", day5_orm_select.DbName)
	//fmt.Println(datasourceName)
	err = orm.RegisterDataBase("default", "mysql", datasourceName)
	if err != nil {
		return err
	}
	return nil
}

func getOrm() orm.Ormer {
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		log.Printf("命令执行失败: %v", err)
		return nil
	}
	return orm.NewOrm()
}
