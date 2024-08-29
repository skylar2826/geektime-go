package day4_session

import (
	__template_and_file "geektime-go/day3_template_and_file"
	"geektime-go/day4_session/session"
	"net/http"
)

var m session.Manager

func RunSession() {

	m.SessionKey = "_sess"
	server := __template_and_file.NewServer("test-session", __template_and_file.ServerWithFilterBuilders(SessionFilter))
	server.Route(http.MethodGet, "/user", handleUser)
	server.Route(http.MethodPost, "/login", handleLogin)
	server.Route(http.MethodPost, '/logout', handleLogout)
}
