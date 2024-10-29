package day9

import (
	"context"
	"math/rand"
	"time"
)

type RandomExpirationCache struct {
	memoryCache
}

func (r *RandomExpirationCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if expiration > 0 {
		offset := time.Duration(rand.Intn(300)) * time.Second // [0, 300)
		expiration = expiration + offset
	}
	return r.Set(ctx, key, value, expiration)
}
