package reflect

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// 元数据注册中心
type Register struct {
	models map[reflect.Type]*Model
}

func NewRegister() *Register {
	return &Register{
		models: make(map[reflect.Type]*Model, 64),
	}
}

func (r *Register) get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)

	m, ok := r.models[typ]
	if !ok {
		var err error
		m, err = r.ParseModel(val)
		if err != nil {
			return nil, err
		}
		// 普通map并发场景下会有读写问题
		r.models[typ] = m
	}
	return m, nil
}

func (r *Register) ParseModel(entity any) (*Model, error) {
	if entity == nil {
		return nil, fmt.Errorf("entity is nil")
	}
	typ := reflect.TypeOf(entity)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()

	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity is not a struct")
	}

	numFields := typ.NumField()

	res := &Model{
		TableName: CamelToSnake(typ.Name()),
		Fields:    make(map[string]*field, numFields),
	}

	for i := 0; i < numFields; i++ {
		fld := typ.Field(i)
		res.Fields[fld.Name] = &field{
			ColName: CamelToSnake(fld.Name),
		}
	}
	return res, nil
}

type Model struct {
	TableName string
	Fields    map[string]*field
}

type field struct {
	ColName string
}

func CamelToSnake(s string) string {
	// str.replace(/([A-Z])/g, "_$1").toLowerCase().slice(1)
	replaced := regexp.MustCompile(`([A-Z])`).ReplaceAllStringFunc(s, func(m string) string {
		return "_" + strings.ToLower(m)
	})

	if len(replaced) != 0 && replaced[0] == '_' {
		replaced = replaced[1:]
	}

	return replaced
}

// double check

// 元数据注册中心
type RegisterV1 struct {
	models map[reflect.Type]*Model
	lock   sync.RWMutex
}

func (r *RegisterV1) getV1(val any) (*Model, error) {
	typ := reflect.TypeOf(val)

	r.lock.RLock()
	m, ok := r.models[typ]
	r.lock.RUnlock()

	if ok {
		return m, nil
	}

	r.lock.Lock()
	m, ok = r.models[typ]
	if ok {
		return m, nil
	}
	var err error
	m, err = r.ParseModel(val)
	if err != nil {
		return nil, err
	}
	r.models[typ] = m

	defer r.lock.Unlock()
	return m, nil
}
func (r *RegisterV1) ParseModel(entity any) (*Model, error) {
	if entity == nil {
		return nil, fmt.Errorf("entity is nil")
	}
	typ := reflect.TypeOf(entity)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()

	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity is not a struct")
	}

	numFields := typ.NumField()

	res := &Model{
		TableName: CamelToSnake(typ.Name()),
		Fields:    make(map[string]*field, numFields),
	}

	for i := 0; i < numFields; i++ {
		fld := typ.Field(i)
		res.Fields[fld.Name] = &field{
			ColName: CamelToSnake(fld.Name),
		}
	}
	return res, nil
}
