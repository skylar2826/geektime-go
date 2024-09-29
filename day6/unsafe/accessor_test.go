package unsafe

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestUser struct {
	Age  int
	name string
}

func TestUnsafeAccessor_Field(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantErr error
		wantRes any
	}{
		{
			name: "field",
			entity: &TestUser{
				Age:  24,
				name: "lily",
			},
			wantRes: 24,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accessor := NewUnsafeAccessor(tc.entity)
			val, err := accessor.Field("Age")
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, val)
		})
	}
}

func TestUnsafeAccessor_SetField(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantErr error
		wantRes any
	}{
		{
			name: "set field",
			entity: &TestUser{
				Age:  24,
				name: "lily",
			},
			wantRes: 24,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accessor := NewUnsafeAccessor(tc.entity)
			err := accessor.SetField("Age", 18)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			var age any
			age, err = accessor.Field("Age")
			require.NoError(t, err)
			assert.Equal(t, age, 18)
		})
	}
}
