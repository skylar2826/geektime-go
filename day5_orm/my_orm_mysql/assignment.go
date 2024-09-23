package my_orm_mysql

type Assignment struct {
	Column string
	Val    any
}

func (Assignment) assign() {}

func Assign(column string, val any) Assignment {
	return Assignment{Column: column, Val: val}
}
