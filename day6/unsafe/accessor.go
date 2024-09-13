package unsafe

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Accessor struct {
	Fields  map[string]*Field
	address unsafe.Pointer
}

func NewUnsafeAccessor(entity any) *Accessor {
	typ := reflect.TypeOf(entity).Elem()
	numField := typ.NumField()

	fields := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		field := typ.Field(i)
		fields[field.Name] = &Field{
			offset: field.Offset,
			typ:    field.Type,
		}
	}

	val := reflect.ValueOf(entity)
	accessor := &Accessor{
		Fields:  fields,
		address: val.UnsafePointer(),
	}

	return accessor
}

func (u *Accessor) Field(name string) (any, error) {
	field, ok := u.Fields[name]
	if !ok {
		return nil, fmt.Errorf("field %s not found", name)
	}

	fieldAddress := unsafe.Pointer(uintptr(u.address) + field.offset)

	return reflect.NewAt(field.typ, fieldAddress).Elem().Interface(), nil
}

func (u *Accessor) SetField(name string, value any) error {
	field, ok := u.Fields[name]

	if !ok {
		return fmt.Errorf("field %s not found", name)
	}

	fieldAddress := unsafe.Pointer(uintptr(u.address) + field.offset)

	reflect.NewAt(field.typ, fieldAddress).Elem().Set(reflect.ValueOf(value))
	return nil
}

type Field struct {
	offset uintptr
	typ    reflect.Type
}
