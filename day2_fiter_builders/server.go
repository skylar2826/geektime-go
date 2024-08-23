package __filter_builder

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Routable interface {
	Route(method string, pattern string, handler func(c *Context)) error
}

type Server interface {
	Routable
	Start(address string) error
	Shutdown(ctx context.Context) error
	Use(method string, pattern string, builders ...FilterBuilder) error
}

type sdkHttpServer struct {
	Name    string
	handler Handler
	root    Filter
	ctxPool sync.Pool
}

func (s *sdkHttpServer) Shutdown(ctx context.Context) error {
	fmt.Printf("正在shutdown server: %s\n", s.Name)
	time.Sleep(2 * time.Second)
	fmt.Printf("shutdown server %s 完成\n", s.Name)
	return nil
}

func (s *sdkHttpServer) Route(method string, pattern string, handler func(c *Context)) error {
	return s.handler.Route(method, pattern, handler)
}

func (s *sdkHttpServer) Start(address string) error {
	//http.Handle("/", s.handler)
	//http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	c := NewContext(writer, request)
	//	s.root(c) // 调用责任链
	//})
	//return http.ListenAndServe(address, nil)
	return http.ListenAndServe(address, s)
}

func (s *sdkHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := s.ctxPool.Get().(*Context)
	defer s.ctxPool.Put(c)
	c.Reset(w, r)
	s.root(c)
}

// use 会执行路由匹配，只有匹配上的mdls才会执行
func (s *sdkHttpServer) Use(method string, pattern string, builders ...FilterBuilder) error {
	return s.handler.AddFilterBuilders(method, pattern, builders...)
}

func NewServer(name string, builders ...FilterBuilder) Server {
	//handler := NewHandlerBasedOnMap()
	handler := NewHandlerBasedOnTree()

	var root Filter = func(c *Context) {
		handler.ServeHTTP(c)
	}
	root = HandleResp(root) // 在最里面就最后执行完成

	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root) // 生成root责任链，由外向里
	}

	return &sdkHttpServer{Name: name, handler: handler, root: root, ctxPool: sync.Pool{
		New: func() interface{} {
			return NewContext()
		},
	}}
}

func NewServerWithFilterNames(name string, filterNames ...string) Server {
	builders := make([]FilterBuilder, 0, len(filterNames))
	for _, name := range filterNames {
		b := GetFilterBuilder(name)
		builders = append(builders, b)
	}

	return NewServer(name, builders...)
}
