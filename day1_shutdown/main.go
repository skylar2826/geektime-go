package __shutdown

import (
	"context"
	"fmt"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func RunShutdown() {

	shutdown := NewGracefulShutdown()
	RegisterFilterBuilder("shutdownFilter", shutdown.ShutdownFilter)
	server := NewServerWithFilterNames("router1", "handleLog", "handleAccess", "shutdownFilter")
	//server := router.NewServer("router1", router.HandleLog, router.HandleAccess, shutdown.ShutdownFilter, router.HandleTracer)
	//server := router.NewServer("router1", router.HandleLog, router.HandleAccess, shutdown.ShutdownFilter, router.HandlePrometheus)
	server.Route(http.MethodPost, "/signUp", SignUp)
	server.Route(http.MethodGet, "/home", Home)
	server.Route(http.MethodGet, "/user/:id", User)
	server.Route(http.MethodGet, "/home/:homeId([0-9][a-zA-Z]+.{3})", Home)

	staticHandler := NewStaticResourceHandler(
		"static",
		"/static",
		WithMoreExtension(map[string]string{
			"mp3": "audio/mp3",
		}),
		WithFileCache(1<<20, 100))
	// 访问 Get http://localhost:8080/static/forest.png
	server.Route(http.MethodGet, "/static/*", staticHandler.ServeStaticResource)

	// adminServer := router.NewServer("router2", router.HandleLog, router.HandleAccess, shutdown.ShutdownFilter)
	adminServer := NewServerWithFilterNames("router2", "handleLog", "handleAccess", "shutdownFilter")
	adminServer.Route(http.MethodGet, "/user", User)

	go func() {
		server.Start("127.0.0.1:8000")
	}()

	go func() {
		adminServer.Start("127.0.0.1:8001")
	}()

	//go func() {
	//	http.Handle("/metrics", promhttp.Handler())
	//	http.ListenAndServe("127.0.0.1:8002", nil)
	//}()

	WaitForShutdown(
		func(ctx context.Context) error {
			fmt.Println("mock notify gateway")
			return nil
		},
		shutdown.RejectNewRequestAndWaiting,
		BuildCloseServerHook(server, adminServer),
		func(ctx context.Context) error {
			fmt.Println("释放资源")
			return nil
		},
	)
}
