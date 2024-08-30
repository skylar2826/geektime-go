package day4_session

import (
	__template_and_file "geektime-go/day3_template_and_file"
	"log"
	"net/http"
)

func SessionFilter(next __template_and_file.Filter) __template_and_file.Filter {
	return func(c *__template_and_file.Context) {
		if c.R.URL.Path == "/login" {
			next(c)
			return
		}
		_, err := m.GetSession(c)
		if err != nil {
			c.W.WriteHeader(http.StatusUnauthorized)
			c.W.Write([]byte("请重新登录"))
			//c.RespUnAuthed()
			return
		}

		// 刷新session的过期时间
		err = m.RefreshSession(c)
		if err != nil {
			log.Printf("session刷新失败: %v\n", err)
		}
		next(c)
	}
}
