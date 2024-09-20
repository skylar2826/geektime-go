package my_orm_mysql

type Assignment struct {
	column string
	val    any
}

func (Assignment) assign() {}

func Assign(column string, val any) Assignment {
	return Assignment{column: column, val: val}
}
