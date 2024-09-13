package valuer

import (
	"database/sql"
	rft "geektime-go/day5_orm/reflect"
	"reflect"
)

type ReflectValue struct {
	val   any
	model *rft.Model
}

var _ Creator = NewReflectValue

func NewReflectValue(model *rft.Model, val any) Valuer {
	return &ReflectValue{val: val, model: model}
}

func (r *ReflectValue) SetColumns(rows *sql.Rows) error {
	cs, err := rows.Columns()
	var vals []any

	for _, c := range cs {
		val := reflect.New(r.model.ColumnMap[c].Typ).Interface()
		vals = append(vals, val)
	}

	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	tpValueElem := reflect.ValueOf(r.val).Elem()
	for i, c := range cs {
		fd := r.model.ColumnMap[c]
		tpValueElem.FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
	}

	return nil
}
