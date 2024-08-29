package day4_session

import (
	"fmt"
	__template_and_file "geektime-go/day3_template_and_file"
)

func handleUser(c *__template_and_file.Context) {
	sess, _ := m.GetSession(c)
	//if err != nil { // 中间件已经判断过了，走到这sess存在
	//	c.RespUnAuthed()
	//}

	// 假设获取nickname字段的值
	val, err := sess.Get(c.R.Context(), "nickname")
	if err != nil {
		c.SystemError()
		return
	}
	c.RespOk()
	fmt.Printf("/user get nickname: %s\n", val)
}

func handleLogin(c *__template_and_file.Context) {
	// 在此之前完成登录注册

	// user信息只能从context里来
	sess, err := m.InitSession(c)
	if err != nil {
		c.SystemError()
		return
	}

	// mock
	err = sess.Set(c.R.Context(), "nickname", "zly")
	if err != nil {
		c.SystemError()
		return
	}
	c.RespOk("登录成功")
}

func handleLogout(c *__template_and_file.Context) {
	err := m.RemoveSession(c)
	if err != nil {
		c.SystemError()
		return
	}
}
