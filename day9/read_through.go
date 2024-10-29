package day9

import (
	"context"
	"golang.org/x/sync/singleflight"
	"time"
)

// ReadThroughCache 是装饰器
// expiration 是过期时间, loadFunc\expiration\logFunc 必填
type ReadThroughCache[T any] struct {
	memoryCache
	loadFunc   func(ctx context.Context, key string) (*T, error) // 加载数据库中的该字段
	expiration time.Duration
	logFunc    func(val ...any)
	g          singleflight.Group
}

func (r *ReadThroughCache[T]) Get(ctx context.Context, key string) (T, error) {

	val, err := r.memoryCache.Get(ctx, key)
	if err != nil && err.Error() == ErrNotFound.Error() {
		val, err = r.loadFunc(ctx, key)
		if err == nil {
			er := r.memoryCache.Set(ctx, key, val, r.expiration)
			if er != nil {
				r.logFunc(er)
			}
		}
	}
	return val.(T), nil
}

func (r *ReadThroughCache[T]) GetV1(ctx context.Context, key string) (T, error) {
	val, err := r.memoryCache.Get(ctx, key)
	if err != nil && err.Error() == ErrNotFound.Error() {
		go func() {
			val, err = r.loadFunc(ctx, key)
			if err == nil {
				er := r.memoryCache.Set(ctx, key, val, r.expiration)
				if er != nil {
					r.logFunc(er)
				}
			}
		}()
	}
	return val.(T), nil
}

func (r *ReadThroughCache[T]) GetV2(ctx context.Context, key string) (T, error) {
	val, err := r.memoryCache.Get(ctx, key)
	if err != nil && err.Error() == ErrNotFound.Error() {
		val, err = r.loadFunc(ctx, key)
		if err == nil {
			go func() {
				er := r.memoryCache.Set(ctx, key, val, r.expiration)
				if er != nil {
					r.logFunc(er)
				}
			}()
		}

	}
	return val.(T), nil
}

func (r *ReadThroughCache[T]) GetV3(ctx context.Context, key string) (T, error) {
	val, err := r.memoryCache.Get(ctx, key)
	if err != nil && err.Error() == ErrNotFound.Error() {
		val, err, _ = r.g.Do(key, func() (interface{}, error) {
			v, er := r.loadFunc(ctx, key)
			if err == nil {
				e := r.memoryCache.Set(ctx, key, v, r.expiration)
				if e != nil {
					r.logFunc(e)
				}
			}
			return v, er
		})

	}
	return val.(T), nil
}

// 使用方式
//type user struct {
//	Name string
//}
//
//func test() {
//	r := &ReadThroughCache[user]{
//		loadFunc: func(ctx context.Context, key string) (user, error) {
//			// todo 从db中获取该字段
//
//		},
//		logFunc: func(err error) {
//			log.Println(err)
//		},
//		expiration: time.Second * 3,
//	}
//	u, err := r.Get(context.Background(), "user_1")
//	u.Name
//}
