//go:build e2e

package day9

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedis_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	rdb.Ping(context.Background())
	c := NewRedisCache(rdb)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := c.Set(ctx, "key1", "value1", time.Minute)
	assert.NoError(t, err)
	var val any
	val, err = c.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, val, "value1")
}

func TestRedis_e2e_V1(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	rdb.Ping(context.Background())
	testCases := []struct {
		name       string
		before     func()
		after      func(t *testing.T)
		wantErr    error
		key        string
		value      any
		expiration time.Duration
	}{
		{
			name:       "set key",
			key:        "key2",
			value:      "value2",
			expiration: time.Second,
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				val, err := rdb.Get(ctx, "key2").Result()
				require.NoError(t, err)
				assert.Equal(t, val, "value2")
				_, err = rdb.Del(ctx, "key2").Result()
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewRedisCache(rdb)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := c.Set(ctx, tc.key, tc.value, tc.expiration)
			assert.Equal(t, err, tc.wantErr)
			tc.after(t)
			if err != nil {
				return
			}
		})
	}
}
