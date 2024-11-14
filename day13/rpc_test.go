package rpc

import (
	"context"
	"errors"
	"geektime-go/day13/user_service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 客户端send 要调用的服务名、方法名和参数
// 服务端找到并调用对应服务和方法，返回处理结果
func TestRPC(t *testing.T) {
	s, err := NewServer("tcp", "127.0.0.1:8081")
	if err != nil {
		t.Log(err)
	}

	u := &user_service.UserServiceServer{}
	s.registerService(u)
	go func() {
		err = s.Start()
		if err != nil {
			t.Log(err)
		}
	}()
	// 等待服务启动
	time.Sleep(time.Second * 6)

	us := &user_service.UserService{}
	err = InitClientProxy("tcp", "127.0.0.1:8081", time.Minute, us)
	if err != nil {
		t.Log(err)
	}

	testCases := []struct {
		name     string
		mock     func()
		wantResp *user_service.GetByIdResponse
		wantErr  error
	}{
		{
			name: "no error",
			mock: func() {
				u.Msg = "hello world"
			},
			wantResp: &user_service.GetByIdResponse{
				Data: "hello world",
			},
		},
		{
			name: "error",
			mock: func() {
				u.Err = errors.New("error")
			},
			wantResp: &user_service.GetByIdResponse{
				Data: "",
			},
			wantErr: errors.New("error"),
		},
		{
			name: "both",
			mock: func() {
				u.Msg = "hello world"
				u.Err = errors.New("error")
			},
			wantErr: errors.New("error"),
			wantResp: &user_service.GetByIdResponse{
				Data: "hello world",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			var res *user_service.GetByIdResponse
			res, err = u.GetById(context.Background(), &user_service.GetByIdRequest{Id: 123})
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, res, tc.wantResp)
			u.Msg = ""
			u.Err = nil
		})
	}

}
