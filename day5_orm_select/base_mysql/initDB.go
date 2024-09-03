package base_mysql

import (
	"database/sql"
	"fmt"
	day5_orm_select "geektime-go/day5_orm_select/types"
)

func InitDB() (*sql.DB, error) {
	datasourceName := fmt.Sprint(day5_orm_select.UserName, ":", day5_orm_select.Password, "@tcp(", day5_orm_select.Ip, ":", day5_orm_select.Port, ")/", day5_orm_select.DbName)

	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
