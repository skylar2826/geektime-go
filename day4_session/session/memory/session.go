package memory

import (
	"context"
	"errors"
	"geektime-go/day4_session/session"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	errorKeyNotFound     = errors.New("session: key not found")
	errorSessionNotFound = errors.New("session: session not found")
)

type Session struct {
	id     string
	values sync.Map
	//mutex  sync.RWMutex
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Get(ctx context.Context, key string) (interface{}, error) {
	val, ok := s.values.Load(key)
	if !ok {
		return nil, errorKeyNotFound
	}
	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, value interface{}) error {
	//s.mutex.RLock()
	//defer s.mutex.RUnlock()
	s.values.Store(key, value)
	return nil
}

type Store struct {
	expiration time.Duration
	sessions   *cache.Cache
	mutex      sync.RWMutex
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		expiration: expiration,
	}
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	sess, ok := s.sessions.Get(id)
	if !ok {
		return nil, errorSessionNotFound
	}
	return sess.(*Session), nil
}

func (s *Store) Generator(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	sess := &Session{
		id: id,
	}
	s.sessions.Set(id, sess, s.expiration)
	return sess, nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.sessions.Delete(id)
	return nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	sess, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	s.sessions.Set(id, sess, s.expiration)
	return nil
}
