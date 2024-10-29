package day9

import (
	"context"
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrFailedToPreemptLock = errors.New("加锁失败")
	ErrLockNotExist        = errors.New("锁不存在")
	//go:embed lua/unlock.lua
	luaUnlock string
)

type lock struct {
	key    string
	value  string
	client redis.Cmdable
}

func (l *lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotExist
	}
	return nil
}

type RedisLockClient struct {
	client redis.Cmdable
}

func NewRedisLockClient(client redis.Cmdable) *RedisLockClient {
	return &RedisLockClient{client}
}

func (r *RedisLockClient) tryLock(ctx context.Context, key string, expiration time.Duration) (*lock, error) {
	val := uuid.New().String()
	ok, err := r.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, ErrFailedToPreemptLock
	}
	if !ok {
		// 代表别人抢到锁
		return nil, ErrFailedToPreemptLock
	}

	return &lock{
		key:    key,
		value:  val,
		client: r.client,
	}, nil
}
