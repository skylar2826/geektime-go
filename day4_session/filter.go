package day4_session

import (
	__template_and_file "geektime-go/day3_template_and_file"
	"log"
)

func SessionFilter(next __template_and_file.Filter) __template_and_file.Filter {
	return func(c *__template_and_file.Context) {
		if c.R.URL.Path == "/login" {
			next(c)
			return
		}
		_, err := m.GetSession(c)
		if err != nil {
			c.RespUnAuthed(err)
			return
		}

		// 刷新session的过期时间
		err = m.RefreshSession(c)
		if err != nil {
			log.Printf("session刷新失败: %v", err)
		}
		next(c)
	}
}
