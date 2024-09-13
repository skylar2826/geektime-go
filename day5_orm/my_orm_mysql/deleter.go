package my_orm_mysql

import (
	rft "geektime-go/day5_orm/model"
	"strings"
)

type Deleter[T any] struct {
	table string
	model *rft.Model

	where []Predicate

	Builder
	r *rft.Register
}

func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}

func (d *Deleter[T]) Where(p ...Predicate) *Deleter[T] {
	if d.where == nil {
		d.where = make([]Predicate, 0, len(p))
	}
	d.where = append(d.where, p...)
	return d
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.sb = &strings.Builder{}
	var err error
	d.model, err = d.r.ParseModel(new(T))
	if err != nil {
		return nil, err
	}
	d.sb.WriteString("delete from ")
	if d.table != "" {
		d.sb.WriteString(d.table)
	} else {
		d.sb.WriteString("`")
		d.sb.WriteString(d.model.TableName)
		d.sb.WriteString("`")
	}

	if len(d.where) > 0 {
		d.sb.WriteString(" where ")
		if err := d.buildPredicate(d.where, d.model); err != nil {
			return nil, err
		}

	}

	d.sb.WriteByte(';')

	return &Query{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}
