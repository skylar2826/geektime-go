package reflect

import (
	"geektime-go/day5_orm/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModels(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantRes *Model
		wantErr error
	}{
		{
			name:   "struct",
			entity: types.TestModel{},
			wantRes: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"LastName": {
						ColName: "last_name",
					},
					"Age": {
						ColName: "age",
					},
				},
			},
		},
		// 用单元测试固化不太合理的测试用例，以免忘记

	}

	for _, tc := range testCases {
		r := NewRegister()
		t.Run(tc.name, func(t *testing.T) {
			model, err := r.ParseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, model)
		})
	}
}

func TestStr(t *testing.T) {
	testCases := []struct {
		name    string
		str     string
		wantRes string
	}{
		{
			name:    "upper cases",
			str:     "ID",
			wantRes: "i_d",
		},
		{
			name:    "use number",
			str:     "Table1Name",
			wantRes: "table1_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := CamelToSnake(tc.str)

			assert.Equal(t, tc.wantRes, s)
		})
	}
}

func TestRegister_get(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantRes   *Model
		wantErr   error
		cacheSize int
	}{
		{
			name:      "cacheSize = 1",
			entity:    types.TestModel{},
			cacheSize: 1,
			wantRes: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"LastName": {
						ColName: "last_name",
					},
					"Age": {
						ColName: "age",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		r := NewRegister()
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, m)
			//assert.Equal(t, tc.cacheSize, len(r.models))
			//
			//typ := reflect.TypeOf(tc.entity)
			//var mo *Model
			//mo, ok := r.models[typ]
			//assert.True(t, ok)
			//assert.Equal(t, mo, tc.wantRes)
		})
	}
}

func TestRegister_tag(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantRes   *Model
		wantErr   error
		cacheSize int
	}{
		{
			name: "tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first_name_t"`
				}
				return &TagTable{}
			}(),

			wantRes: &Model{
				TableName: "tag_table",
				FieldMap: map[string]*Field{

					"FirstName": {
						ColName: "first_name_t",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		r := NewRegister()
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, m)
		})
	}
}

type customUserName1 struct {
	FirstName string
}

type customUserName struct {
	FirstName string
}

func (c customUserName) TableName() string {
	return "custom_table"
}

func TestRegister_tableName(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantRes   *Model
		wantErr   error
		cacheSize int
	}{
		{
			name:   "no custom table name",
			entity: customUserName1{},
			wantRes: &Model{
				TableName: "custom_user_name1",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
					},
				},
			},
		},
		{
			name:   "custom table name",
			entity: customUserName{},
			wantRes: &Model{
				TableName: "custom_table",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
					},
				},
			},
		},
		{
			name:   "custom table name ptr",
			entity: &customUserName{},
			wantRes: &Model{
				TableName: "custom_table",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		r := NewRegister()
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, m)
		})
	}
}
