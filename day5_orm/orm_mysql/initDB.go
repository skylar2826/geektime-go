package orm_mysql

import (
	"fmt"
	day5_orm_select "geektime-go/day5_orm/types"
	"github.com/astaxie/beego/orm"
	"log"
)

func initDB() error {
	orm.RegisterModel(new(day5_orm_select.User))
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		return err
	}
	datasourceName := fmt.Sprint(day5_orm_select.UserName, ":", day5_orm_select.Password, "@tcp(", day5_orm_select.Ip, ":", day5_orm_select.Port, ")/", day5_orm_select.DbName)
	fmt.Println(datasourceName)
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
