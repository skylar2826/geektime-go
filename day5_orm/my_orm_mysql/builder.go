package my_orm_mysql

import (
	"fmt"
	rft "geektime-go/day5_orm/model"
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
		if exp.op != "" {
			s.sb.WriteString(" " + exp.op.String() + " ")
		}
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
		err := s.buildColumn(exp, model)
		if err != nil {
			return err
		}
	case value:
		s.sb.WriteString("?")
		s.addArgs(exp.val)
	case RawExpr:
		s.sb.WriteString("(")
		s.sb.WriteString(exp.expression)
		s.addArgs(exp.args...)
		s.sb.WriteString(")")
	default:
		return fmt.Errorf("invalid expression type: %T", expr)
	}
	return nil
}

func (s *Builder) buildColumn(col Column, model *rft.Model) error {
	s.sb.WriteByte('`')
	name, ok := model.FieldMap[col.name]
	if !ok {
		return fmt.Errorf("field %s not found", col.name)
	}
	s.sb.WriteString(name.ColName)
	s.sb.WriteByte('`')
	s.buildAs(col)

	return nil
}

func (s *Builder) buildAs(col Column) {
	if col.alias != "" {
		s.sb.WriteString(" AS ")
		s.sb.WriteString("`")
		s.sb.WriteString(col.alias)
		s.sb.WriteString("`")
	}
}

func (s *Builder) addArgs(val ...any) *Builder {
	if len(val) == 0 {
		return s
	}
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, val...)
	return s
}

func (s *Builder) BuildColumns(columns []Selectable, model *rft.Model) error {
	if len(columns) == 0 {
		s.sb.WriteString("*")
		return nil
	}

	for i, c := range columns {
		if i > 0 {
			s.sb.WriteString(",")
		}
		switch col := c.(type) {
		case Column:
			err := s.buildColumn(c.(Column), model)
			if err != nil {
				return err
			}
		case Aggregate:
			// 聚合函数名
			s.sb.WriteString(col.fn)
			s.sb.WriteString("(`")
			s.sb.WriteString(col.arg)

			s.sb.WriteString("`)")
			s.buildAs(Column{name: col.arg, alias: col.alias})
		case RawExpr:
			s.sb.WriteString(col.expression)
			s.addArgs(col.args...)
		}

	}
	return nil
}
