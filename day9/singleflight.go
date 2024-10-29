package day9

import (
	"context"
	"golang.org/x/sync/singleflight"
	"log"
	"time"
)

type singleFlightCache[T any] struct {
	*ReadThroughCache[T]
}

func NewSingleFlightCache[T any](cache memoryCache, loadFunc func(cxt context.Context, key string) (any, error), expiration time.Duration) *singleFlightCache[T] {
	g := singleflight.Group{}
	return &singleFlightCache[T]{
		ReadThroughCache: &ReadThroughCache[T]{
			memoryCache: cache,
			logFunc:     log.Fatal,
			expiration:  expiration,
			loadFunc: func(ctx context.Context, key string) (*T, error) {
				val, err, _ := g.Do(key, func() (interface{}, error) {
					return loadFunc(ctx, key)
				})
				v := val.(T)
				return &v, err
			},
		},
	}
}
