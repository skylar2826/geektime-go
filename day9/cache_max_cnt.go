package day9

import (
	"context"
	"errors"
	"time"
)

var (
	ErrOutMaxCnt = "超出最大限制"
)

type CacheMaxCnt struct {
	*BuildInMapCache
	maxCnt int32
	cnt    int32
}

func NewCacheMaxCnt(maxCnt int32) *CacheMaxCnt {
	b := NewBuildInMapCache()

	res := &CacheMaxCnt{
		BuildInMapCache: b,
		maxCnt:          maxCnt,
		cnt:             int32(0),
	}

	origin := res.onEvicted

	res.onEvicted = func(key string, val any) {
		res.cnt--
		origin(key, val)
	}
	return res
}

func (c *CacheMaxCnt) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	c.BuildInMapCache.mutex.Lock()
	_, ok := c.BuildInMapCache.data[key]
	if !ok {
		if c.cnt+1 > c.maxCnt {
			return errors.New(ErrOutMaxCnt)
		}
		c.cnt++
		err := c.BuildInMapCache.set(ctx, key, value, expiration)
		return err
	}
	return c.BuildInMapCache.set(ctx, key, value, expiration)
}
