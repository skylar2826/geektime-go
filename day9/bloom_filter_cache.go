package day9

import (
	"context"
	"log"
	"time"
)

type BloomFilter struct {
	HasKey func(ctx context.Context, key string) bool
}

type BloomFilterCache[T any] struct {
	*ReadThroughCache[T]
	bf BloomFilter
}

func NewBloomFilterCache[T any](cache memoryCache, bf BloomFilter, loadFunc func(cxt context.Context, key string) (*T, error), expiration time.Duration) *BloomFilterCache[T] {
	return &BloomFilterCache[T]{
		ReadThroughCache: &ReadThroughCache[T]{
			memoryCache: cache,
			loadFunc: func(ctx context.Context, key string) (*T, error) {
				if bf.HasKey(ctx, key) {
					return loadFunc(ctx, key)
				}
				return nil, ErrNotFound
			},
			expiration: expiration,
			logFunc:    log.Fatal,
		},
		bf: bf,
	}
}

type BloomFilterCacheV1[T any] struct {
	*ReadThroughCache[T]
	bf BloomFilter
}

func (b *BloomFilterCacheV1[T]) Get(ctx context.Context, key string) (T, error) {
	val, err := b.memoryCache.Get(ctx, key)
	if err != nil && err.Error() == ErrNotFound.Error() {
		if b.bf.HasKey(ctx, key) {
			var v any
			v, err = b.loadFunc(ctx, key)
			if err != nil {
				er := b.memoryCache.Set(ctx, key, v, b.expiration)
				if er != nil {
					b.logFunc(er)
				}
			}
		}
	}
	return val.(T), nil

}
