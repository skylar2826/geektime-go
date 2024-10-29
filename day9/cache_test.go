package day9

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		cache   func() *BuildInMapCache
		wantVal any
		wantErr error
	}{
		{
			name: "not key",
			key:  "not key",
			cache: func() *BuildInMapCache {
				return NewBuildInMapCache()
			},
			wantErr: ErrNotFound,
		},
		{
			name: "exist key",
			key:  "key1",
			cache: func() *BuildInMapCache {
				b := NewBuildInMapCache()
				err := b.Set(context.Background(), "key1", 123, time.Second)
				require.NoError(t, err)
				return b

			},
			wantVal: 123,
		},
		{
			name: "time out",
			key:  "time out key",
			cache: func() *BuildInMapCache {
				b := NewBuildInMapCache()
				err := b.Set(context.Background(), "time out key", 123, time.Second)
				require.NoError(t, err)
				time.Sleep(time.Second * 2)
				return b
			},
			wantErr: errors.New(ErrNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.cache().Get(context.Background(), tc.key)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, val, tc.wantVal)
		})
	}
}

func TestBuildInMapCache_Loop(t *testing.T) {
	cnt := 0
	b := NewBuildInMapCache(func(b *BuildInMapCache) {
		cnt++
	})

	err := b.Set(context.Background(), "key1", 123, time.Second)
	require.NoError(t, err)

	time.Sleep(time.Second * 3)

	b.mutex.RLock()
	defer b.mutex.RUnlock()
	_, ok := b.data["key1"]
	require.False(t, ok)
	require.Equal(t, cnt, 1)
}
