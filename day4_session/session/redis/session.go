package redis

import (
	"context"
	"errors"
	"fmt"
	"geektime-go/day4_session/session"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	errSessionNotFound = errors.New("session not found")
)

type Session struct {
	redisKey string
	id       string
	client   redis.Cmdable
}

func (s *Session) Set(ctx context.Context, key string, value interface{}) error {
	// lua 下标从1开始 lua 解决线程
	const lua = `
	if redis.call("exists", KEYS[1)
	then
		return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
	else
		return -1
	end
	`
	res, err := s.client.Eval(ctx, lua, []string{s.redisKey}, key, value).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errSessionNotFound
	}
	return nil
}

func (s *Session) Get(ctx context.Context, key string) (interface{}, error) {
	/*
		如果要检查session是否过期，需要在此处使用lua脚本检查
	*/
	val, err := s.client.HGet(ctx, redisKey, key).Result() // 没判断是否过期；如果过期了，就get不到; 用户使用差异不大
	return val, err
}

func (s *Session) ID() string {
	return s.id
}

type Store struct {
	client     redis.Cmdable
	expiration time.Duration
	prefix     string
}

type StoreOption func(*Store)

func NewStore(client redis.Cmdable, opts ...StoreOption) *Store {
	res := &Store{
		expiration: time.Minute * 15,
		prefix:     "sessionId",
		client:     client,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func StoreWithPrefix(prefix string) StoreOption {
	return func(s *Store) {
		s.prefix = prefix
	}
}

func StoreWithExpiration(expiration time.Duration) StoreOption {
	return func(s *Store) {
		s.expiration = expiration
	}
}

func (s *Store) Generator(ctx context.Context, id string) (session.Session, error) {
	// 不需要考虑线程安全，id全局唯一，不存在重复使用

	key := setRedisKey(s.prefix, id)

	//	const lua = `
	//redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
	//return redis.call("expire", KEYS[1], ARGV[3])
	//`
	//
	//	_, err := s.client.Eval(ctx, lua, []string{key}, "_sess_id", id, s.expiration.Milliseconds()).Result()
	//	if err != nil {
	//		return nil, err
	//	}

	// id, id 占位， 值无意义；后面会用真实的 sessionId={key, value}中的key, value覆盖
	_, err := s.client.HSet(ctx, key, id, id).Result()
	if err != nil {
		return nil, err
	}
	_, err = s.client.Expire(ctx, id, s.expiration).Result()
	if err != nil {
		return nil, err
	}

	return &Session{
		id:       id,
		redisKey: key,
		client:   s.client,
	}, nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	key := setRedisKey(s.prefix, id)
	_, err := s.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
	// 代表id不存在
	//if cnt == 0 {
	//	return errSessionNotFound
	//}
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	key := setRedisKey(s.prefix, id)
	ok, err := s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("session: id 对应的 session 不存在")
	}
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	/*
		自由决策要不要提前把session存储的用户数据一并捞过来
		1. 都不拿 （此处使用）
		2. 只拿热点数据/高频数据
		3. 都拿
	*/

	/*
		get的时候要不要判断存不存在？
		可以不判断，等到真正拿的时候在判断；提前判断也没用，可能刚判断存在，就过期了；返回的&session{}再拿的时候就不存在；所以没有必要
	*/
	key := setRedisKey(s.prefix, id)

	cnt, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, errSessionNotFound
	}

	return &Session{
		id:       id,
		redisKey: key,
		client:   s.client,
	}, nil
}

func setRedisKey(prefix string, id string) string {
	return fmt.Sprintf("%s-%s", prefix, id)
}
