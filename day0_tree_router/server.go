package __tree_router

import "net/http"

type Routable interface {
	Route(method string, pattern string, handler func(c *Context)) error
}

type Server interface {
	Routable
	Start(address string) error
}

type sdkHttpServer struct {
	Name    string
	handler Handler
	root    Filter
}

func (s *sdkHttpServer) Route(method string, pattern string, handler func(c *Context)) error {
	return s.handler.Route(method, pattern, handler)
}

func (s *sdkHttpServer) Start(address string) error {
	//http.Handle("/", s.handler)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		c := NewContext(writer, request)
		s.root(c) // 调用责任链
	})
	return http.ListenAndServe(address, nil)
}

func NewServer(name string, builders ...FilterBuilder) Server {
	//handler := NewHandlerBasedOnMap()
	handler := NewHandlerBasedOnTree()

	var root Filter = func(c *Context) {
		handler.ServeHTTP(c)
	}

	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root) // 生成root责任链，由外向里
	}

	return &sdkHttpServer{Name: name, handler: handler, root: root}
}
