package my_orm_mysql

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type selector[T any] struct {
	table string
	where []Predicate
	sb    *strings.Builder
	args  []any
}

func (s *selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	s.sb.WriteString("select * from ")

	if s.table != "" {
		s.sb.WriteString(s.table)
	} else {
		var t T
		typ := reflect.TypeOf(t)
		s.sb.WriteByte('`')
		s.sb.WriteString(typ.Name())
		s.sb.WriteByte('`')
	}

	if len(s.where) > 0 {
		s.sb.WriteString(" where ")
		p := s.where[0]
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}

		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
			if err := s.buildExpression(p); err != nil {
				return nil, err
			}
		}

	}

	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *selector[T]) Get(ctx context.Context) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *selector[T]) From(table string) *selector[T] {
	s.table = table
	return s
}

func (s *selector[T]) Where(f ...Predicate) *selector[T] {
	if s.where == nil {
		s.where = make([]Predicate, 0, len(f))
	}
	s.where = append(s.where, f...)
	return s
}

func (s *selector[T]) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil: // 因为有default throw error, 所以Not左边没有是nil需要用case处理
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
		s.sb.WriteString(" " + exp.op.String() + " ")
		_, ok = exp.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(exp.name)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteString("?")
		s.addArgs(exp.val)

	default:
		return fmt.Errorf("invalid expression type: %T", expr)
	}
	return nil
}

func (s *selector[T]) addArgs(val any) *selector[T] {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, val)
	return s
}
