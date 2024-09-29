package my_orm_mysql

import (
	"context"
	"database/sql"
	"geektime-go/day5_orm/internal"
	"strings"
)

// Selectable 是一个标记接口
// 它代表要查找的列或者聚合方法
type Selectable interface {
	selectable()
}

type Selector[T any] struct {
	table tableReference
	where []Predicate

	//r *rft.Register
	//db       *DB
	columns  []Selectable
	groupBy  []Column
	having   []Predicate
	orderBys []OrderBy
	offset   int
	limit    int

	sess Session
	Builder
}

func NewSelector[T any](sess Session) *Selector[T] {
	builder := NewBuilder(sess)
	return &Selector[T]{
		//db:      db,
		Builder: *builder,
		sess:    sess,
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var err error
	if s.model == nil {
		s.model, err = s.R.ParseModel(new(T))
		if err != nil {
			return nil, err
		}
	}
	s.sb.WriteString("select ")
	if err = s.BuildColumns(s.columns); err != nil {
		return nil, err
	}
	s.sb.WriteString(" from ")

	if err = s.buildTable(s.table); err != nil {
		return nil, err
	}

	if len(s.where) > 0 {
		s.sb.WriteString(" where ")

		if err := s.buildPredicate(s.where); err != nil {
			return nil, err
		}

	}

	if len(s.groupBy) > 0 {
		s.sb.WriteString(" group by ")
		for idx, col := range s.groupBy {
			if idx > 0 {
				s.sb.WriteString(",")
			}
			err = s.buildColumn(col)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(s.having) > 0 {
		s.sb.WriteString(" having ")
		if err := s.buildPredicate(s.having); err != nil {
			return nil, err
		}
	}

	if len(s.orderBys) > 0 {
		s.sb.WriteString(" order by ")
		for idx, ob := range s.orderBys {
			if idx > 0 {
				s.sb.WriteString(",")
			}
			err = s.buildColumn(ob.col)
			if err != nil {
				return nil, err
			}
			s.sb.WriteString(" " + ob.order)
		}
	}

	// limit在offset前
	if s.limit != 0 {
		s.sb.WriteString(" limit ?")
		s.addArgs(s.limit)
	}

	if s.offset != 0 {
		s.sb.WriteString(" offset ?")
		s.addArgs(s.offset)
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
	rows, err = s.sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, internal.ErrorNoRows
	}

	tp := new(T)
	v := s.Creator(s.model, tp)
	//var valuer valuer2.Valuer
	err = v.SetColumns(rows)
	return tp, err
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	var err error
	s.model, err = s.R.ParseModel(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, s.sess, s.core, &QueryContext{
		Type:    "Select",
		Builder: s,
		Model:   s.model,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	var err error
	s.model, err = s.R.ParseModel(new(T))
	if err != nil {
		return nil, err
	}
	return getMuti[T](ctx, s.sess, s.core, &QueryContext{
		Type:    "Select",
		Builder: s,
		Model:   s.model,
	})
}

func (s *Selector[T]) From(table tableReference) *Selector[T] {
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

func (s *Selector[T]) GroupBy(col ...Column) *Selector[T] {
	s.groupBy = col
	return s
}

func (s *Selector[T]) Having(p ...Predicate) *Selector[T] {
	s.having = p
	return s
}

func (s *Selector[T]) orderBy(orderBys ...OrderBy) *Selector[T] {
	s.orderBys = orderBys
	return s
}
func (s *Selector[T]) Offset(offset int) *Selector[T] {
	s.offset = offset
	return s
}
func (s *Selector[T]) Limit(limit int) *Selector[T] {
	s.limit = limit
	return s
}
