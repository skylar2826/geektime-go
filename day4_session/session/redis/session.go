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
	errorKeyNotFound     = errors.New("session: key not found")
	errorSessionNotFound = errors.New("session: session not found")
)

type Session struct {
	id       string
	redisKey string
	client   redis.Cmdable
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := s.client.HGet(ctx, s.redisKey, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, value interface{}) error {
	const lua = `
if redis.call("exists", KEYS[1])
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
		return fmt.Errorf("session: session创建失败")
	}
	return nil
}

type Store struct {
	client     redis.Cmdable
	expiration time.Duration
	prefix     string
}

// 防止和redis里的其他key重名
func getRedisKey(prefix string, key string) string {
	return fmt.Sprintf("%s-%s", prefix, key)
}

type StoreOption func(s *Store)

func StoreWithExpiration(expiration time.Duration) StoreOption {
	return func(s *Store) {
		s.expiration = expiration
	}
}

func StoreWithPrefix(prefix string) StoreOption {
	return func(s *Store) {
		s.prefix = prefix
	}
}

func NewStore(client redis.Cmdable, opts ...StoreOption) *Store {
	res := &Store{client: client}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	key := getRedisKey(s.prefix, id)

	cnt, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	// 不存在
	if cnt == 0 {
		return nil, errorSessionNotFound
	}
	return &Session{
		id:       id,
		redisKey: key,
		client:   s.client,
	}, nil
}

func (s *Store) Generator(ctx context.Context, id string) (session.Session, error) {
	key := getRedisKey(s.prefix, id)

	_, err := s.client.HSet(ctx, key, "", "").Result()
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
	key := getRedisKey(s.prefix, id)

	_, err := s.client.HDel(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	key := getRedisKey(s.prefix, id)
	ok, err := s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errorSessionNotFound
	}
	return nil
}
