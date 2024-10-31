package day9

import (
	"context"
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"time"
)

var (
	ErrFailedToPreemptLock = errors.New("加锁失败")
	ErrLockNotExist        = errors.New("锁不存在")
	ErrFailedLock          = errors.New("过期时间更新失败")
	//go:embed lua/unlock.lua
	luaUnlock string
	//go:embed lua/refresh.lua
	luaRefresh string
)

type Lock struct {
	key        string
	value      string
	client     redis.Cmdable
	expiration time.Duration
	closeChan  chan struct{}
}

func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value).Result()
	l.closeChan <- struct{}{}
	if err != nil {
		return err
	}
	if res != int64(1) {
		return ErrLockNotExist
	}
	return nil
}

func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.value, l.expiration.Seconds()).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrFailedLock
	}
	return nil
}

func (l *Lock) AutoRefresh(expiration time.Duration, interval time.Duration) error {
	timeoutChan := make(chan struct{}, 1)
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), expiration)
			err := l.Refresh(ctx)
			cancel()
			if errors.Is(err, context.DeadlineExceeded) {
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {
				close(timeoutChan)
				return err
			}

		case <-l.closeChan:
			close(timeoutChan)
			return nil
		case <-timeoutChan:

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := l.Refresh(ctx)
			cancel()
			if errors.Is(err, context.DeadlineExceeded) {
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {
				close(timeoutChan)
				return err
			}
		}
	}
}

type RedisLockClient struct {
	client redis.Cmdable
	g      singleflight.Group
}

func NewRedisLockClient(client redis.Cmdable) *RedisLockClient {
	return &RedisLockClient{client: client}
}

// SingleFlightLock 处理高并发场景
func (r *RedisLockClient) SingleFlightLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	for {
		flag := false
		resCh := r.g.DoChan(key, func() (interface{}, error) {
			flag = true
			return r.tryLock(ctx, key, expiration)
		})
		select {
		case res := <-resCh:
			if flag { //该实例抢到锁
				r.g.Forget(key)
				if res.Err != nil {
					return res.Val.(*Lock), nil
				}
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (r *RedisLockClient) tryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	val := uuid.New().String()
	ok, err := r.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, ErrFailedToPreemptLock
	}
	if !ok {
		// 代表别人抢到锁
		return nil, ErrFailedToPreemptLock
	}

	return &Lock{
		key:        key,
		value:      val,
		client:     r.client,
		expiration: expiration,
	}, nil
}
