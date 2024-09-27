package my_orm_mysql

import (
	"geektime-go/day5_orm/internal/valuer"

	"geektime-go/day5_orm/model"
)

type core struct {
	model       *model.Model
	dialect     Dialect
	Creator     valuer.Creator
	R           *model.Register
	middlewares []Middleware
}
