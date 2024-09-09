package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTraverse(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantErr error
		wantRes []any
	}{
		{
			name:    "array",
			entity:  [3]int{1, 2, 3},
			wantRes: []any{1, 2, 3},
		},
		{
			name:    "slice",
			entity:  []int{1, 2, 3},
			wantRes: []any{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateArrayOrSlice(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestTraverseMap(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantErr error
		wantRes map[string]any
	}{
		{
			name: "map",
			entity: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			wantRes: map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateMap(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
