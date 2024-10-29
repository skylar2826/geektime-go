package day9

import (
	"context"
	"time"
)

type WriteThroughCache struct {
	memoryCache
	storeFunc  func(ctx context.Context, key string, value any) error
	logFunc    func(e error)
	expiration time.Duration
}

func (w *WriteThroughCache) Set(ctx context.Context, key string, value any) error {
	err := w.storeFunc(ctx, key, value)
	if err != nil {
		return err
	}
	return w.memoryCache.Set(ctx, key, value, w.expiration)
}

func (w *WriteThroughCache) SetV1(ctx context.Context, key string, value any) error {
	err := w.storeFunc(ctx, key, value)
	if err != nil {
		return err
	}
	go func() {
		er := w.memoryCache.Set(ctx, key, value, w.expiration)
		if er != nil {
			w.logFunc(er)
		}
	}()
	return nil
}

type WriteThroughCacheV1[T any] struct {
	memoryCache
	storeFunc  func(ctx context.Context, key string, value any) error
	logFunc    func(e error)
	expiration time.Duration
}

func (w *WriteThroughCacheV1[T]) Set(ctx context.Context, key string, value any) error {
	err := w.storeFunc(ctx, key, value.(T))
	if err != nil {
		return err
	}
	return w.memoryCache.Set(ctx, key, value, w.expiration)
}
