package my_orm_mysql

import (
	"context"
	"database/sql"
	"geektime-go/day5_orm/internal"
	model2 "geektime-go/day5_orm/model"
	rft "geektime-go/day5_orm/reflect"
	"strings"
)

// Selectable 是一个标记接口
// 它代表要查找的列或者聚合方法
type Selectable interface {
	selectable()
}

type Selector[T any] struct {
	table string
	where []Predicate
	model *model2.Model
	Builder
	//r *rft.Register
	db      *rft.DB
	columns []Selectable
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
	s.sb.WriteString("select ")
	if err = s.BuildColumns(s.columns, s.model); err != nil {
		return nil, err
	}
	s.sb.WriteString(" from ")

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

func (s *Selector[T]) GetV1(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = s.db.DB.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, internal.ErrorNoRows
	}

	tp := new(T)
	v := s.db.Creator(s.model, tp)
	//var valuer valuer2.Valuer
	err = v.SetColumns(rows)
	return tp, err
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = s.db.DB.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, internal.ErrorNoRows
	}

	tp := new(T)
	v := s.db.Creator(s.model, tp)
	err = v.SetColumns(rows)
	if err != nil {
		return nil, err
	}
	return tp, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = s.db.DB.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	var tps []*T
	for rows.Next() {
		tp := new(T)
		v := s.db.Creator(s.model, tp)
		err = v.SetColumns(rows)
		if err != nil {
			return nil, err
		}

		tps = append(tps, tp)
	}
	return tps, nil
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

func (s *Selector[T]) Select(field ...Selectable) *Selector[T] {
	s.columns = field
	return s
}
