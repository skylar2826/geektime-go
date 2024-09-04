package orm_mysql

import (
	day5_orm_select "geektime-go/day5_orm/types"
	"github.com/astaxie/beego/orm"
)

func insert(o orm.Ormer) error {
	user := new(day5_orm_select.User)
	user.Name = "xiaozhu"
	_, err := o.Insert(user)
	if err != nil {
		return err
	}
	return nil
}
