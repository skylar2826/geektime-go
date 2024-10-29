package day9

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrFailedKey = "字段设置失败"
)

type RedisCache struct {
	c redis.Cmdable
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	res, err := r.c.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return err
	}
	if res != "OK" {
		return errors.New(ErrFailedKey)
	}
	return nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, error) {
	return r.c.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	_, err := r.c.Del(ctx, key).Result()
	return err
}

func NewRedisCache(ctrl redis.Cmdable) *RedisCache {
	return &RedisCache{
		c: ctrl,
	}
}
