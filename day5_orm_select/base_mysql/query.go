package base_mysql

import (
	"database/sql"
	day5_orm_select "geektime-go/day5_orm_select/types"
)

func Query(db *sql.DB) ([]day5_orm_select.User, error) {
	var user day5_orm_select.User
	var users []day5_orm_select.User

	// 使用prepare复用语句
	//stmt, err := db.Prepare("select * from user where id = ?")
	//if err != nil {
	//	return nil, err
	//}
	//defer stmt.Close()
	//
	//var rows *sql.Rows
	//rows, err = stmt.Query(3)

	rows, err := db.Query("select * from user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex, &user.Phone)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
