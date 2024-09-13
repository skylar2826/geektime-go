package unsafe

import (
	"fmt"
	"reflect"
	"unsafe"
)

type UnsafeAccessor struct {
	Fields  map[string]*FieldMeta
	address unsafe.Pointer
}

func NewUnsafeAccessor(entity any) *UnsafeAccessor {
	typ := reflect.TypeOf(entity).Elem()

	numField := typ.NumField()
	fields := make(map[string]*FieldMeta, numField)

	val := reflect.ValueOf(entity)

	for i := 0; i < numField; i++ {
		field := typ.Field(i)
		fields[field.Name] = &FieldMeta{
			offset: field.Offset,
			typ:    field.Type,
		}
	}

	return &UnsafeAccessor{
		Fields: fields,
		// 	val.UnsafeAddr() 返回uintptr是个值，GC（垃圾回收)后，address改变，uintptr失效; val.UnsafePointer 返回 unsafe.pointer是一个指针，GC会维护更新它的值
		address: val.UnsafePointer(),
	}
}

func (u *UnsafeAccessor) Field(name string) (any, error) {
	field, ok := u.Fields[name]
	if !ok {
		return nil, fmt.Errorf("field %s not found", name)
	}

	fieldAddress := unsafe.Pointer(uintptr(u.address) + field.offset)

	// 知道地址 =》 把地址转化成*对象
	return reflect.NewAt(field.typ, fieldAddress).Elem().Interface(), nil
	// 已知返回类型
	//return *(*int)(fieldAddress), nil
}

func (u *UnsafeAccessor) SetField(name string, value any) error {
	field, ok := u.Fields[name]
	if !ok {
		return fmt.Errorf("field %s not found", name)
	}
	fieldAddress := unsafe.Pointer(uintptr(u.address) + field.offset)
	reflect.NewAt(field.typ, fieldAddress).Elem().Set(reflect.ValueOf(value))
	return nil
}

type FieldMeta struct {
	//Offset unsafe.Pointer
	offset uintptr
	typ    reflect.Type
}
