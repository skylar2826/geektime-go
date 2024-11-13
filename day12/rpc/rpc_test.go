package rpc

import (
	"context"
	"geektime-go/day12/rpc/user_service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 客户端send 要调用的服务名、方法名和参数
// 服务端找到并调用对应服务和方法，返回处理结果
func TestRPC(t *testing.T) {
	go func() {
		s, err := NewServer("tcp", "127.0.0.1:8081")
		if err != nil {
			t.Log(err)
		}

		u := &user_service.UserServiceServer{}
		s.registerService(u)

		err = s.Start()
		if err != nil {
			t.Log(err)
		}
	}()
	// 等待服务启动
	time.Sleep(time.Second * 6)

	u := &user_service.UserService{}
	err := InitClientProxy("tcp", "127.0.0.1:8081", time.Minute, u)
	if err != nil {
		t.Log(err)
	}
	var res *user_service.GetByIdResponse
	res, err = u.GetById(context.Background(), &user_service.GetByIdRequest{Id: 123})
	if err != nil {
		t.Log(err)
	}
	assert.Equal(t, res, &user_service.GetByIdResponse{
		Data: "请求Id: 123, 响应信息：hello world",
	})
}
