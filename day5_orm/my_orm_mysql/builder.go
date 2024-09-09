package my_orm_mysql

import (
	"fmt"
	rft "geektime-go/day5_orm/reflect"
	"strings"
)

type Builder struct {
	sb   *strings.Builder
	args []any
}

func NewBuilder() *Builder {
	return &Builder{
		sb: &strings.Builder{},
	}
}

func (s *Builder) buildPredicate(pd []Predicate, model *rft.Model) error {
	p := pd[0]
	if err := s.buildExpression(p, model); err != nil {
		return err
	}

	for i := 1; i < len(pd); i++ {
		p = p.And(pd[i])
		if err := s.buildExpression(p, model); err != nil {
			return err
		}
	}
	return nil
}

func (s *Builder) buildExpression(expr Expression, model *rft.Model) error {
	switch exp := expr.(type) {
	case nil: // 因为有default throw error, 所以Not左边没有是nil需要用case处理
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left, model); err != nil {
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
		if err := s.buildExpression(exp.right, model); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
	case Column:
		s.sb.WriteByte('`')
		name, ok := model.Fields[exp.name]
		if !ok {
			return fmt.Errorf("field %s not found", exp.name)
		}
		s.sb.WriteString(name.ColName)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteString("?")
		s.addArgs(exp.val)

	default:
		return fmt.Errorf("invalid expression type: %T", expr)
	}
	return nil
}

func (s *Builder) addArgs(val any) *Builder {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, val)
	return s
}
