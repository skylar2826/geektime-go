package net

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNet(t *testing.T) {
	go func() {
		s, err := NewServer("tcp", "127.0.0.1:8081")
		if err != nil {
			t.Log(err)
		}

		err = s.Start()
		if err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second * 5)

	for i := 0; i < 3; i++ {
		// 每次需要创建新的连接，不能复用之前的conn? 为啥？因为之前的conn不知道回没回来（被释放）？
		client, err := NewClient("tcp", "127.0.0.1:8081", time.Minute)
		if err != nil {
			t.Log(err)
		}
		var res []byte
		res, err = client.Send(context.Background(), []byte(fmt.Sprintf("request %d", i)))
		if err != nil {
			t.Log(err)
		}
		fmt.Println("响应：", string(res))

		//// 测试连接可以被复用，等待链接释放
		//time.Sleep(time.Second * 7)
	}
}
