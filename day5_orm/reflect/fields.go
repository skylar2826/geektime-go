package reflect

import (
	internal "geektime-go/day5_orm/internal"
	"reflect"
)

func IterateFields(entity any) (map[string]any, error) {
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if val.IsZero() {
		return nil, internal.ErrorValueIsZero
	}

	for typ.Kind() == reflect.Ptr {
		// *User
		typ = typ.Elem()
		// (*User)(nil)
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, internal.ErrorEntityNotStruct
	}

	numField := typ.NumField()

	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		fld := typ.Field(i)
		if fld.IsExported() {
			res[fld.Name] = val.Field(i).Interface()
		} else {
			res[fld.Name] = reflect.Zero(fld.Type).Interface()
		}
	}

	return res, nil
}

func SetField(entity any, field string, value any) error {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	val = val.FieldByName(field)
	if !val.CanSet() {
		return internal.ErrorFieldCantSet
	}
	val.Set(reflect.ValueOf(value))

	return nil
}
