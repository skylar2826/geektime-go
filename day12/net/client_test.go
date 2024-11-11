package net

import (
	"fmt"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	go func() {
		//err := Serve("tcp", "127.0.0.1:8082")
		s := &Server{}
		err := s.Start("tcp", "127.0.0.1:8082")
		t.Log(err)
	}()

	time.Sleep(time.Second * 5)

	//err := Connect("tcp", "127.0.0.1:8082", time.Minute)
	c, err := NewClient("tcp", "127.0.0.1:8082", time.Minute)
	t.Log(err)
	var repData string
	repData, err = c.Send("hello world")
	t.Log(err)
	fmt.Println("响应：", repData)

}
