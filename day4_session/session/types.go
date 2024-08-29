package session

import (
	"context"
	"net/http"
)

// session 直接与内存|redis交互， 存储key|value键值对，有sessionId
// 相当于二维Map  map[string]map[string]any   { sessionId: { key: value } }
type Session interface {
	ID() string // id是私有的，外部通过ID()获取
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}) error
}

type Store interface {
	Get(ctx context.Context, id string) (Session, error)
	Generator(ctx context.Context, id string) (Session, error)
	Remove(ctx context.Context, id string) error
	Refresh(ctx context.Context, id string) error
}

type Propagator interface {
	Extract(r *http.Request) (string, error)
	Inject(id string, w http.ResponseWriter) error
	Remove(w http.ResponseWriter) error
}
