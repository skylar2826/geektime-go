package model

import (
	"geektime-go/day5_orm/internal"
	"geektime-go/day5_orm/types"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var (
	tagName = "column"
)

type Register struct {
	Models sync.Map
}

func NewRegister() *Register {
	return &Register{
		//models: make(map[reflect.Type]*Model, 64),
	}
}

type Model struct {
	TableName string
	FieldMap  map[string]*Field
	ColumnMap map[string]*Field
}

type Field struct {
	ColName string
	Typ     reflect.Type
	GoName  string
	Offset  uintptr
}

func (r *Register) ParseModel(entity any) (*Model, error) {
	if entity == nil {
		return nil, internal.ErrorEntityIsNil
	}
	typ := reflect.TypeOf(entity)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, internal.ErrorEntityNotStruct
	}

	var tableName string
	if tbl, ok := entity.(types.TableName); ok {
		tableName = tbl.TableName()
	}
	if tableName == "" {
		tableName = CamelToSnake(typ.Name())
	}

	numFields := typ.NumField()
	FieldMap := make(map[string]*Field, numFields)
	ColumnMap := make(map[string]*Field, numFields)
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)

		tagPair, err := r.parseTag(field.Tag)
		if err != nil {
			return nil, err
		}
		fieldName := tagPair[tagName]
		if fieldName == "" {
			fieldName = CamelToSnake(field.Name)
		}

		FieldMap[field.Name] = &Field{
			ColName: fieldName,
			Typ:     field.Type,
			GoName:  field.Name,
			Offset:  field.Offset,
		}
		ColumnMap[fieldName] = FieldMap[field.Name]
	}

	m := &Model{
		TableName: tableName,
		FieldMap:  FieldMap,
		ColumnMap: ColumnMap,
	}

	return m, nil
}

type User struct {
	Id string `orm:"column=uid;xxx=x"`
}

func (r *Register) parseTag(tag reflect.StructTag) (map[string]string, error) {
	val, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}

	res := map[string]string{}
	pairs := strings.Split(val, ";")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return map[string]string{}, internal.ErrorTagFormat
		}
		res[kv[0]] = kv[1]
	}
	return res, nil
}

func (r *Register) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)

	//model, ok := r.models[typ]
	model, ok := r.Models.Load(typ)
	if !ok {
		var err error
		model, err = r.ParseModel(val)
		if err != nil {
			return nil, err
		}
		//r.models[typ] = model
		r.Models.Store(typ, model)
	}

	return model.(*Model), nil
}

// CamelToSnake FirstName => first_name
func CamelToSnake(s string) string {
	res := regexp.MustCompile(`([A-Z])`).ReplaceAllStringFunc(s, func(m string) string {
		return "_" + strings.ToLower(m)
	})

	if len(res) != 0 && res[0] == '_' {
		res = res[1:]
	}
	return res
}
