package __filter_builder

import (
	"net/http"
)

func RunFilterBuilders() {
	server := NewServer("router1", HandleLog, HandleAccess)

	// 注册路由
	server.Route(http.MethodGet, "/a/b/d/e", Home)

	// 注册middleware
	server.Use(http.MethodGet, "/", HandlerX)
	server.Use(http.MethodGet, "/a", HandlerA)
	server.Use(http.MethodGet, "/a/b", HandlerAB)
	server.Use(http.MethodGet, "/a/d", HandlerAD)
	server.Use(http.MethodGet, "/a/*", HandlerAX)
	server.Use(http.MethodGet, "/a/b/d", HandlerABD)
	server.Use(http.MethodGet, "/a/b/d/e", HandlerABDE)
	server.Use(http.MethodGet, "/a/b/d/*", HandlerABDX)

	server.Start("127.0.0.1:8000")
}
