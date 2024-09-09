package my_orm_mysql

import (
	"context"
	rft "geektime-go/day5_orm/reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	where []Predicate
	model *rft.Model
	Builder
	//r *rft.Register
	db *rft.DB
}

func NewSelector[T any](db *rft.DB) *Selector[T] {
	builder := NewBuilder()
	return &Selector[T]{
		db:      db,
		Builder: *builder,
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var err error
	s.model, err = s.db.R.ParseModel(new(T))
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("select * from ")

	if s.table != "" {
		s.sb.WriteString(s.table)
	} else {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.TableName)
		s.sb.WriteByte('`')
	}

	if len(s.where) > 0 {
		s.sb.WriteString(" where ")

		if err := s.buildPredicate(s.where, s.model); err != nil {
			return nil, err
		}

	}

	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) Get(ctx context.Context) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(f ...Predicate) *Selector[T] {
	if s.where == nil {
		s.where = make([]Predicate, 0, len(f))
	}
	s.where = append(s.where, f...)
	return s
}
