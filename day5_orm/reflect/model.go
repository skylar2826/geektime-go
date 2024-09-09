package reflect

import (
	"fmt"
	"geektime-go/day5_orm/types"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

const (
	tagName = "column"
)

// 元数据注册中心
type Register struct {
	models sync.Map
}

func NewRegister() *Register {
	return &Register{
		//models: make(map[reflect.Type]*Model, 64),
	}
}

func (r *Register) get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)

	m, ok := r.models.Load(typ.String())
	//m, ok := r.models[typ]
	if !ok {
		var err error
		// 多个goroutine 都会进来会有第二个goroutine覆盖前一个的问题，会重复解析和store
		// 刚启动的使用有轻微的覆盖问题，但map的性能比double-check好
		m, err = r.ParseModel(val)
		if err != nil {
			return nil, err
		}
		// 普通map并发场景下会有读写问题
		//r.models[typ] = m
		r.models.Store(typ.String(), m)
	}
	return m.(*Model), nil
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

	var tableName string
	// 有类型断言的地方，指针和struct的测试用例都要写
	if tbl, ok := entity.(types.TableName); ok {
		tableName = tbl.TableName()
	}
	if tableName == "" {
		tableName = CamelToSnake(typ.Name())
	}

	numFields := typ.NumField()

	res := &Model{
		TableName: tableName,
		Fields:    make(map[string]*field, numFields),
	}

	for i := 0; i < numFields; i++ {

		fld := typ.Field(i)
		pair, err := r.ParseTag(fld.Tag)
		if err != nil {
			return nil, err
		}
		columnName := pair[tagName]
		if columnName == "" {
			columnName = CamelToSnake(fld.Name)
		}
		res.Fields[fld.Name] = &field{
			ColName: columnName,
		}
	}
	return res, nil
}

//
//type User struct {
//	id int `orm:"column=id;xxx=x"`
//}

//column=id => {column: id}

func (r *Register) ParseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag := tag.Get("orm")
	if ormTag == "" {
		return map[string]string{}, nil
	}

	res := make(map[string]string, 1)
	pairs := strings.Split(ormTag, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("orm tag format error")
		}
		res[kv[0]] = kv[1]
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
