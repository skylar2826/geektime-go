package base_mysql

import (
	"fmt"
	day5_orm_select "geektime-go/day5_orm_select/types"
	"log"
)

func RunBase() {
	db, err := InitDB()
	if err != nil {
		log.Printf("connect database error:%v", err)
		return
	}
	var users []day5_orm_select.User
	users, err = Query(db)
	if err != nil {
		log.Printf("query user error:%v", err)
	}

	for _, user := range users {
		fmt.Printf("user: %#v\n", user)
	}

}
