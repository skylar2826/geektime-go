package orm_mysql

import (
	"fmt"
	"log"
)

func RunOrm() {
	err := initDB()
	if err != nil {
		log.Fatal("初始化数据库失败:", err)
		return
	}
	o := getOrm()
	err = insert(o)
	if err != nil {
		log.Printf("插入失败: %v", err)
	}
	fmt.Println("插入成功")
}
