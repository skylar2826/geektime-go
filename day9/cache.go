package day9

import (
	"context"
	"errors"
	_ "github.com/go-redis/redis"
	_ "github.com/golang/mock/mockgen/model"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("field not found")
)

type memoryCache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Delete(ctx context.Context, key string) error
}

type item struct {
	val      any
	deadline time.Time
}

func (i *item) deadlineBefore(t time.Time) bool {
	return !i.deadline.IsZero() && i.deadline.Before(t)
}

type BuildInMapCache struct {
	mutex     sync.RWMutex
	data      map[string]*item
	close     chan struct{}
	onEvicted func(key string, val any)
}

type BuildInMapCacheOption func(cache *BuildInMapCache)

func BuildInMapCacheWithEvictedCallback(fn func(key string, val any)) BuildInMapCacheOption {
	return func(cache *BuildInMapCache) {
		cache.onEvicted = fn
	}
}

func NewBuildInMapCache(opts ...BuildInMapCacheOption) *BuildInMapCache {
	res := &BuildInMapCache{
		data:  make(map[string]*item, 100),
		close: make(chan struct{}),
		onEvicted: func(key string, val any) {

		},
	}

	for _, opt := range opts {
		opt(res)
	}

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case t := <-ticker.C:
				res.mutex.Lock()
				i := 0
				for key, value := range res.data {
					if i > 1000 {
						break
					}
					if value.deadlineBefore(t) {
						res.delete(context.Background(), key)
					}
					i++
				}
				res.mutex.Unlock()
			case <-res.close:
				return
			}
		}

	}()

	return res
}

func (b *BuildInMapCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	//if expiration > 0 {
	//	time.AfterFunc(expiration, func() {
	//		b.mutex.Lock()
	//		defer b.mutex.Unlock()
	//		val, ok := b.data[key]
	//		if ok && !val.deadline.IsZero() && val.deadline.Before(time.Now()) {
	//			delete(b.data, key)
	//		}
	//	})
	//}
	return b.set(ctx, key, value, expiration)
}

func (b *BuildInMapCache) set(ctx context.Context, key string, val any, expiration time.Duration) error {
	var dl time.Time
	if expiration > 0 {
		dl = time.Now().Add(expiration)
	}
	b.data[key] = &item{
		val:      val,
		deadline: dl,
	}
	return nil
}

func (b *BuildInMapCache) Get(ctx context.Context, key string) (any, error) {
	b.mutex.RLock()
	res, ok := b.data[key]
	b.mutex.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	t := time.Now()
	if res.deadlineBefore(t) {
		b.mutex.Lock() // double check原因：lock锁住后，可能有别人set值
		defer b.mutex.Unlock()
		res, ok = b.data[key]
		if !ok {
			return nil, ErrNotFound
		}
		if res.deadlineBefore(t) {
			b.delete(ctx, key)
			return nil, ErrNotFound
		}
	}

	return res.val, nil
}

func (b *BuildInMapCache) Close() error {
	b.close <- struct{}{}
	return nil
}

func (b *BuildInMapCache) Delete(ctx context.Context, key string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.delete(ctx, key)
}

func (b *BuildInMapCache) delete(ctx context.Context, key string) {
	res, ok := b.data[key]
	if !ok {
		return
	}
	delete(b.data, key)
	b.onEvicted(key, res)
}
