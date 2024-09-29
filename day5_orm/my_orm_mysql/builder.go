package my_orm_mysql

import (
	"errors"
	"fmt"
	rft "geektime-go/day5_orm/model"
	"strings"
)

type Builder struct {
	sb     *strings.Builder
	args   []any
	quoter byte
	core
}

func (s *Builder) quote(name string) {
	s.sb.WriteByte(s.quoter)
	s.sb.WriteString(name)
	s.sb.WriteByte(s.quoter)
}

func NewBuilder(sess Session) *Builder {
	core := sess.getCore()
	return &Builder{
		sb:     &strings.Builder{},
		quoter: core.dialect.quoter(),
		core:   core,
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
		if err := s.buildColumn(exp); err != nil {
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
	case Aggregate:
		if err := s.BuildAggregate(exp); err != nil {
			return nil
		}
	default:
		return fmt.Errorf("invalid expression type: %T", expr)
	}
	return nil
}

func (s *Builder) buildColumn(col Column) error {
	var field *rft.Field
	var ok bool

	switch t := col.table.(type) {
	case nil:
		field, ok = s.model.FieldMap[col.Name]
		if !ok {
			return fmt.Errorf("field %s not found", col.Name)
		}
	case Table:
		model, err := s.R.Get(t.entity)
		if err != nil {
			return err
		}

		field, ok = model.FieldMap[col.Name]
		if !ok {
			return fmt.Errorf("field %s not found", col.Name)
		}

		if t.alias != "" {
			s.quote(t.alias)
			s.sb.WriteString(".")
		}
	default:
		return errors.New("invalid column type")
	}

	s.quote(field.ColName)
	s.buildAs(col)

	return nil
}

func (s *Builder) buildTable(table tableReference) error {
	switch t := table.(type) {
	case nil:
		s.quote(s.model.TableName)
	case Table:
		m, err := s.R.ParseModel(t.entity)
		if err != nil {
			return err
		}
		s.quote(m.TableName)
		if t.alias != "" {
			s.sb.WriteString(" As ")
			s.quote(t.alias)
		}
	case Join:
		s.sb.WriteString("(")
		if err := s.buildTable(t.left); err != nil {
			return err
		}
		s.sb.WriteString(" " + t.typ + " ")
		if err := s.buildTable(t.right); err != nil {
			return err
		}

		if len(t.using) > 0 {
			s.sb.WriteString(" Using (")
			for idx, using := range t.using {
				if idx > 0 {
					s.sb.WriteString(",")
				}
				if err := s.buildColumn(Column{Name: using}); err != nil {
					return err
				}
			}
			s.sb.WriteString(")")
		}
		if len(t.on) > 0 {
			s.sb.WriteString(" On (")
			if err := s.buildPredicate(t.on); err != nil {
				return err
			}
			s.sb.WriteString(")")
		}

		s.sb.WriteString(")")
	default:
		return errors.New("invalid table type")
	}
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

func (s *Builder) BuildAggregate(col Aggregate) error {
	// 聚合函数名
	s.sb.WriteString(col.fn)
	s.sb.WriteString("(`")
	s.sb.WriteString(col.arg)

	s.sb.WriteString("`)")
	s.buildAs(Column{Name: col.arg, alias: col.alias})
	return nil
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
			if err := s.buildColumn(c.(Column)); err != nil {
				return err
			}
		case Aggregate:
			if err := s.BuildAggregate(col); err != nil {
				return err
			}
		case RawExpr:
			s.sb.WriteString(col.expression)
			s.addArgs(col.args...)
		}

	}
	return nil
}
