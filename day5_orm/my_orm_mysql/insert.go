package my_orm_mysql

import (
	"fmt"
	"geektime-go/day5_orm/internal"
	"geektime-go/day5_orm/model"
	rft "geektime-go/day5_orm/reflect"
	"reflect"
	"strings"
)

// Assignable 标记接口 用于 update 和 upset中
type Assignable interface {
	assign()
}

type OnDuplicateKey[T any] struct {
	assigns []Assignable
}

type OnDuplicateKeyBuilder[T any] struct {
	i *Insert[T]
}

func (o *OnDuplicateKeyBuilder[T]) Update(assigns ...Assignable) *Insert[T] {
	o.i.onDuplicateKey = &OnDuplicateKey[T]{
		assigns: assigns,
	}
	return o.i
}

type Insert[T any] struct {
	db      *rft.DB
	sb      strings.Builder
	values  []*T
	columns []string
	//onDuplicateKey []Assignable
	onDuplicateKey *OnDuplicateKey[T]
}

//func (i *Insert[T]) OnDuplicateKey(assigns ...Assignable) *Insert[T] {
//	i.onDuplicateKey = assigns
//	return i
//}

func (i *Insert[T]) OnDuplicateKey() *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
		i: i,
	}
}

func NewInsert[T any](db *rft.DB) *Insert[T] {
	return &Insert[T]{
		db: db,
	}
}

func (i *Insert[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, internal.ErrorInsertZeroRow
	}
	i.sb.WriteString("INSERT INTO ")

	m, err := i.db.R.Get(i.values[0]) // new(T)
	if err != nil {
		return nil, err
	}

	i.sb.WriteByte('`')
	i.sb.WriteString(m.TableName)
	i.sb.WriteByte('`')

	fields := m.Fields
	if len(i.columns) > 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, goName := range i.columns {
			field, ok := m.FieldMap[goName]
			if !ok {
				return nil, internal.ErrorFieldUnknown
			}
			fields = append(fields, field)
		}
	}

	i.sb.WriteString("(")
	for idx, field := range fields {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('`')
		i.sb.WriteString(field.ColName)
		i.sb.WriteByte('`')
	}
	i.sb.WriteString(")")

	i.sb.WriteString(" VALUES ")

	args := make([]interface{}, 0, len(i.values)*len(fields))
	for idx, row := range i.values {
		if idx > 0 {
			i.sb.WriteByte(',')
		}

		i.sb.WriteString("(")

		for j, field := range fields {
			if j > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteString("?")
			arg := reflect.ValueOf(row).Elem().FieldByName(field.GoName).Interface()
			args = append(args, arg)
		}
		i.sb.WriteString(")")

	}

	if i.onDuplicateKey != nil {
		i.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
		for idx, assign := range i.onDuplicateKey.assigns {
			if idx > 0 {
				i.sb.WriteByte(',')
			}

			switch a := assign.(type) {
			case Assignment:
				i.sb.WriteString("`")
				colName := m.FieldMap[a.column].ColName
				i.sb.WriteString(colName)
				i.sb.WriteString("`")
				i.sb.WriteString("=?")
				args = append(args, a.val)
			case Column:
				i.sb.WriteString("`")
				colName := m.FieldMap[a.name].ColName
				i.sb.WriteString(colName)
				i.sb.WriteString("`")
				i.sb.WriteString("=VALUES(")
				i.sb.WriteString("`")
				i.sb.WriteString(colName)
				i.sb.WriteString("`")
				i.sb.WriteString(")")
			default:
				return nil, fmt.Errorf("未知类型： %T", a)
			}
		}
	}

	i.sb.WriteString(";")

	return &Query{
		SQL:  i.sb.String(),
		Args: args,
	}, nil
}

func (i *Insert[T]) Values(vals ...*T) *Insert[T] {
	i.values = vals
	return i
}

func (i *Insert[T]) Columns(cols ...string) *Insert[T] {
	i.columns = cols
	return i
}
