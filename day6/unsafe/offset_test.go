package unsafe

import "testing"

type User struct {
	Name    string
	Age     int32
	Alias   []string
	Address string
}
type User2 struct {
	Name    string
	Age     int32 // 4字节
	Age2    int32
	Alias   []string
	Address string
}

func TestUnsafe(t *testing.T) {
	testCases := []struct {
		name   string
		entity any
	}{
		{
			name:   "print offset",
			entity: User{},
		},
		{
			name:   "print offset2",
			entity: User2{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			PrintUnsafeOffset(tc.entity)
		})
	}
}
