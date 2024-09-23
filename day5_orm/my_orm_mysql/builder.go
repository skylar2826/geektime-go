package my_orm_mysql

import (
	"fmt"
	rft "geektime-go/day5_orm/model"
	"strings"
)

type Builder struct {
	sb      *strings.Builder
	args    []any
	model   *rft.Model
	dialect Dialect
	quoter  byte
}

func (s *Builder) quote(name string) {
	s.sb.WriteByte(s.quoter)
	s.sb.WriteString(name)
	s.sb.WriteByte(s.quoter)
}

func NewBuilder(db *DB) *Builder {
	return &Builder{
		sb:      &strings.Builder{},
		dialect: db.Dialect,
		quoter:  db.Dialect.quoter(),
	}
}

func (s *Builder) buildPredicate(pd []Predicate) error {
	p := pd[0]
	if err := s.buildExpression(p); err != nil {
		return err
	}

	for i := 1; i < len(pd); i++ {
		p = p.And(pd[i])
		if err := s.buildExpression(p); err != nil {
			return err
		}
	}
	return nil
}

func (s *Builder) buildExpression(expr Expression) error {
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
		if exp.op != "" {
			s.sb.WriteString(" " + exp.op.String() + " ")
		}
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
		err := s.buildColumn(exp)
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

func (s *Builder) buildColumn(col Column) error {

	name, ok := s.model.FieldMap[col.Name]
	if !ok {
		return fmt.Errorf("field %s not found", col.Name)
	}

	s.quote(name.ColName)
	s.buildAs(col)

	return nil
}

func (s *Builder) buildAs(col Column) {
	if col.alias != "" {
		s.sb.WriteString(" AS ")
		s.quote(col.alias)
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

func (s *Builder) BuildColumns(columns []Selectable) error {
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
			err := s.buildColumn(c.(Column))
			if err != nil {
				return err
			}
		case Aggregate:
			// 聚合函数名
			s.sb.WriteString(col.fn)
			s.sb.WriteString("(`")
			s.sb.WriteString(col.arg)

			s.sb.WriteString("`)")
			s.buildAs(Column{Name: col.arg, alias: col.alias})
		case RawExpr:
			s.sb.WriteString(col.expression)
			s.addArgs(col.args...)
		}

	}
	return nil
}
