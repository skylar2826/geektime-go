package reflect

import (
	"fmt"
	"geektime-go/day5_orm/internal"
	types "geektime-go/day5_orm/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestUser struct {
	Id    int
	name  string
	age   int
	Phone string
}

func TestIterateFields(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantRes map[string]any
		wantErr error
	}{
		{
			name: "struct",
			entity: TestUser{
				Id:    1,
				name:  "lily",
				age:   18,
				Phone: "123456",
			},
			wantRes: map[string]any{
				"Id":    1,
				"name":  "", // 拿不到私有值
				"age":   0,
				"Phone": "123456",
			},
		},
		{
			name: "struct",
			entity: &TestUser{
				Id:    1,
				name:  "lily",
				age:   18,
				Phone: "123456",
			},
			wantRes: map[string]any{
				"Id":    1,
				"name":  "", // 拿不到私有值
				"age":   0,
				"Phone": "123456",
			},
		},
		{
			name:    "err kind",
			entity:  18,
			wantErr: internal.ErrorEntityNotStruct,
		},
		{
			name: "mutil ptr",
			entity: func() **TestUser {
				res := &TestUser{
					Id:    1,
					name:  "lily",
					age:   18,
					Phone: "123456",
				}
				return &res
			}(),
			wantRes: map[string]any{
				"Id":    1,
				"name":  "", // 拿不到私有值
				"age":   0,
				"Phone": "123456",
			},
		},
		{
			name:    "nil kind",
			entity:  nil,
			wantErr: internal.ErrorEntityNotStruct,
		},
		{
			name:    "(*user)(nil)",
			entity:  (*TestUser)(nil),
			wantErr: fmt.Errorf("value is zero"),
		},
		{
			name: "测试其他包的对象反射能不能拿到",
			entity: &types.TestPerson{
				Id:   1,
				Name: "lily",
			},
			wantRes: map[string]any{
				"Id":   1,
				"Name": "lily",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entity, err := IterateFields(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, entity)
		})
	}
}

func TestSetField(t *testing.T) {
	testCases := []struct {
		name       string
		entity     any
		field      string
		value      any
		wantErr    error
		wantEntity any
	}{
		{
			name: "unexported field",
			entity: &TestUser{
				Id:   1,
				name: "lily",
			},
			field:   "name",
			value:   "zhuzhu",
			wantErr: internal.ErrorFieldCantSet,
		},
		{
			name: "struct field",
			entity: TestUser{
				Id:   1,
				name: "lily",
			},
			field:   "Id",
			value:   2,
			wantErr: internal.ErrorFieldCantSet,
		},
		{
			name: "struct",
			entity: &TestUser{
				Id:   1,
				name: "lily",
			},
			field: "Id",
			value: 2,

			wantEntity: &TestUser{
				Id:   2,
				name: "lily",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.value)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}
