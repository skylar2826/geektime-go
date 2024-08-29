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

		sess, err := m.GetSession(c)
		if err != nil {
			c.RespUnAuthed()
			return
		}

		// 刷新session
		err = m.Refresh(c.R.Context(), sess.ID())
		if err != nil {
			log.Println("刷新session失败", err)
		}

		next(c)
	}
}
