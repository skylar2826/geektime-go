package valuer

import (
	"database/sql"
	rft "geektime-go/day5_orm/model"
	"reflect"
)

type ReflectValue struct {
	val   reflect.Value
	model *rft.Model
}

var _ Creator = NewReflectValue

func NewReflectValue(model *rft.Model, val any) Valuer {
	return &ReflectValue{val: reflect.ValueOf(val).Elem(), model: model}
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

	tpValueElem := r.val
	for i, c := range cs {
		fd := r.model.ColumnMap[c]
		tpValueElem.FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
	}

	return nil
}

func (r *ReflectValue) Field(name string) (any, error) {
	return r.val.FieldByName(name).Interface(), nil
}
