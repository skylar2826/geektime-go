package day4_session

import (
	__template_and_file "geektime-go/day3_template_and_file"
	"geektime-go/day4_session/session"
	"geektime-go/day4_session/session/cookie"
	memory "geektime-go/day4_session/session/memory"
	"net/http"
	"time"
)

var m *session.Manager

func RunSession() {
	m = &session.Manager{
		Propagator:    cookie.NewPropagator(),
		Store:         memory.NewStore(time.Second * 15),
		CtxSessionKey: "_session",
	}
	server := __template_and_file.NewServer("test-session", __template_and_file.ServerWithFilterBuilders(SessionFilter))
	server.Route(http.MethodPost, "/login", handleLogin)
	server.Route(http.MethodGet, "/user", handleUser)
	server.Route(http.MethodPost, "/logout", handleLogout)
	server.Start("127.0.0.1:8000")
}
