package memory

// 基于内存的实现
import (
	"context"
	"errors"
	"fmt"
	"geektime-go/day4_session/session"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	errorKeyNotFound     = errors.New("key not found")
	errorSessionNotFound = errors.New("session not found")
)

type Store struct {
	sessions   *cache.Cache
	expiration time.Duration
	mutex      sync.RWMutex
}

func NewStore(expiration time.Duration) *Store {
	// 过期时间即cache的过期时间
	return &Store{
		// 1分钟检查一次缓存是否过期
		sessions:   cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

type Session struct {
	values sync.Map
	id     string

	// 方案二: 控制性更强
	//mutex sync.RWMutex // 锁
	//values map[string]any
}

// 默认实现都加上error和context, 习惯；虽然当前不用，但面对第三方实现可能是要用的，context重要

func (s *Session) Set(ctx context.Context, key string, value interface{}) error {
	s.values.Store(key, value)
	return nil
}

func (s *Session) Get(ctx context.Context, key string) (interface{}, error) {
	val, ok := s.values.Load(key)
	if !ok {
		//return nil, fmt.Errorf("%w, key: %s", errorKeyNotFound, key)
		return nil, errorKeyNotFound
	}
	return val, nil
}

func (s *Session) ID() string {
	// 写一次便只读，所以不加锁
	return s.id
}

func (s *Store) Generator(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess := &Session{
		id: id,
	}
	s.sessions.Set(id, sess, s.expiration)
	return sess, nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// 是不存在还是删除成功无法判断
	s.sessions.Delete(id)
	return nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// 有写操作，不是并发安全的，需要加锁
	// cache上没有刷新操作，所以需要重新设置覆盖
	val, ok := s.sessions.Get(id)
	if !ok {
		return fmt.Errorf("session: 该 id 对应的session 不存在, id: %s", id)
	}
	s.sessions.Set(id, val, s.expiration)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	sess, ok := s.sessions.Get(id)
	if !ok {
		return nil, errorSessionNotFound
	}
	return sess.(*Session), nil
}
