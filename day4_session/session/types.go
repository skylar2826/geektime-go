package session

import (
	"context"
	"net/http"
)

// Store 管理session的创建、销毁
type Store interface {
	Generator(ctx context.Context, id string) (Session, error)
	Remove(ctx context.Context, id string) error
	Refresh(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (Session, error)
}

// Session 基于内存的Session
type Session interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	ID() string
}

// Propagator 操作http
type Propagator interface {
	// Inject 将session插入response
	Inject(id string, w http.ResponseWriter) error
	// Extract 在请求中获取sessionId
	Extract(r *http.Request) (string, error)
	// Remove 在response中删除session
	Remove(w http.ResponseWriter) error
}
