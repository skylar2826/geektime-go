package reflect

import (
	"fmt"
	"reflect"
)

func IterateFields(entity any) (map[string]any, error) {
	// typ是nil,但val不是nil
	if entity == nil {
		return nil, fmt.Errorf("entity is nil")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if val.IsZero() { // 对象占据的内存空间都是0
		return nil, fmt.Errorf("不支持零值")
	}

	// 如果是if仅支持*User, for支持多级指针 **User
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity kind must be struct, kind: %v", typ.Kind())
	}

	num := typ.NumField()
	res := make(map[string]any, num)

	for i := 0; i < num; i++ {
		field := typ.Field(i)      // 字段信息
		fieldValue := val.Field(i) // 值信息
		if field.IsExported() {
			res[field.Name] = fieldValue.Interface()
		} else {
			res[field.Name] = reflect.Zero(field.Type).Interface()
		}
	}
	return res, nil
}

func SetField(entity any, field string, value any) error {
	val := reflect.ValueOf(entity)
	// val.type.Kind
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	val = val.FieldByName(field)
	if !val.CanSet() {
		return fmt.Errorf("不能修改")
	}
	val.Set(reflect.ValueOf(value))
	return nil
}
