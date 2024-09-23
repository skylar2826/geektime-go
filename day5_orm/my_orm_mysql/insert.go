package my_orm_mysql

import (
	"geektime-go/day5_orm/internal"
	"geektime-go/day5_orm/model"
	"reflect"
)

// Assignable 标记接口 用于 update 和 upset中
type Assignable interface {
	assign()
}

type Upsert struct {
	Assigns         []Assignable
	conflictColumns []string
}

type UpsertBuilder[T any] struct {
	i               *Insert[T]
	conflictColumns []string
}

func (o *UpsertBuilder[T]) Update(assigns ...Assignable) *Insert[T] {
	o.i.upsert = &Upsert{
		Assigns:         assigns,
		conflictColumns: o.conflictColumns,
	}
	return o.i
}

// ConflictColumns 中间方法
func (o *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	o.conflictColumns = cols
	return o
}

type Insert[T any] struct {
	db      *DB
	values  []*T
	columns []string
	upsert  *Upsert
	Builder
}

func (i *Insert[T]) Upsert() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}

func NewInsert[T any](db *DB) *Insert[T] {
	builder := NewBuilder(db)
	return &Insert[T]{
		db:      db,
		Builder: *builder,
	}
}

func (i *Insert[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, internal.ErrorInsertZeroRow
	}
	var err error
	i.model, err = i.db.R.ParseModel(new(T))
	i.sb.WriteString("INSERT INTO ")

	m, err := i.db.R.Get(i.values[0]) // new(T)
	if err != nil {
		return nil, err
	}

	i.quote(m.TableName)

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
		i.quote(field.ColName)
	}
	i.sb.WriteString(")")

	i.sb.WriteString(" VALUES ")

	i.args = make([]interface{}, 0, len(i.values)*len(fields))
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
			i.addArgs(arg)
		}
		i.sb.WriteString(")")

	}

	if i.upsert != nil {
		err = i.dialect.upsert(&i.Builder, i.upsert)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteString(";")

	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
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
