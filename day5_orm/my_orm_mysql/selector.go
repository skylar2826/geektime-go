package my_orm_mysql

import (
	"context"
	"database/sql"
	"fmt"
	"geektime-go/day5_orm/internal"
	rft "geektime-go/day5_orm/reflect"
	"reflect"
	"strings"
	"unsafe"
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
	address := reflect.ValueOf(tp).UnsafePointer()
	var cs []string
	cs, err = rows.Columns()
	var vals []any
	if err != nil {
		return nil, err
	}

	for _, c := range cs {
		fd, ok := s.model.ColumnMap[c]
		if !ok {
			return nil, fmt.Errorf("column %s not found", c)
		}
		fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
		val := reflect.NewAt(fd.Typ, fdAddress).Interface()
		vals = append(vals, val)
	}

	err = rows.Scan(vals...)
	if err != nil {
		return nil, err
	}

	return tp, nil
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

	var cs []string
	cs, err = rows.Columns()
	var vals []any

	for _, c := range cs {
		val := reflect.New(s.model.ColumnMap[c].Typ).Interface()
		vals = append(vals, val)
	}

	err = rows.Scan(vals...)
	if err != nil {
		return nil, err
	}

	tp := new(T)
	tpValueElem := reflect.ValueOf(tp).Elem()
	for i, c := range cs {
		fd := s.model.ColumnMap[c]
		tpValueElem.FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
	}

	return tp, err
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//q, err := s.Build()
	//if err != nil {
	//	return nil, err
	//}
	//
	//var rows *sql.Rows
	//rows, err = s.db.DB.QueryContext(ctx, q.SQL, q.Args)
	//if err != nil {
	//	return nil, err
	//}
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
