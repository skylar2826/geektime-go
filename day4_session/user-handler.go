package day4_session

import (
	"fmt"
	__template_and_file "geektime-go/day3_template_and_file"
	"log"
)

func handleLogin(c *__template_and_file.Context) {
	// 在InitSession之前校验用户名、密码
	session, err := m.InitSession(c)
	if err != nil {
		c.SystemError(err)
		return
	}
	err = session.Set(c.R.Context(), "nickname", "zly")
	if err != nil {
		c.SystemError(err)
		return
	}
	c.RespOk("登录成功")
}

func handleUser(c *__template_and_file.Context) {
	session, _ := m.GetSession(c) // 已经经过中间件sessionFilter，所以一定能拿到session, 可以忽略错误
	//if err != nil {
	//	c.RespUnAuthed()
	//	return
	//}
	// 假设要把昵称从session中拿出
	sessionValue, err := session.Get(c.R.Context(), "nickname")
	if err != nil {
		log.Printf("获取昵称失败 %s\n", err)
		return
	}
	fmt.Printf("/user 获取昵称：%s\n", sessionValue)
	c.RespOk()
}

func handleLogout(c *__template_and_file.Context) {
	err := m.RemoveSession(c)
	if err != nil {
		c.SystemError(err)
		return
	}
	c.RespOk("退出成功")
}
