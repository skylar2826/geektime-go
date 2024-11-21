package rpc

import (
	"context"
	"errors"
	"geektime-go/day13/serialize/proto"
	"geektime-go/day13/serialize/proto/gen"
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
	var c *Client
	c, err = NewClient("tcp", "127.0.0.1:8081", time.Minute)
	if err != nil {
		t.Log(err)
	}
	err = c.InitService(us)
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
			res, err = us.GetById(context.Background(), &user_service.GetByIdRequest{Id: 123})
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, res, tc.wantResp)
			u.Msg = ""
			u.Err = nil
		})
	}

}

func TestRPCProto(t *testing.T) {
	s, err := NewServer("tcp", "127.0.0.1:8081")
	if err != nil {
		t.Log(err)
	}

	u := &user_service.UserServiceServer{}
	s.registerService(u)
	s.registerSerialize(&proto.Serializer{})

	go func() {
		err = s.Start()
		if err != nil {
			t.Log(err)
		}
	}()
	// 等待服务启动
	time.Sleep(time.Second * 6)

	us := &user_service.UserService{}
	var c *Client
	c, err = NewClient("tcp", "127.0.0.1:8081", time.Minute, ClientWithSerialize(&proto.Serializer{}))
	if err != nil {
		t.Log(err)
	}
	err = c.InitService(us)
	if err != nil {
		t.Log(err)
	}

	testCases := []struct {
		name    string
		Id      int64
		wantErr error
	}{

		{
			name: "normal",
			Id:   int64(123456),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			var res *gen.GetByIdResp
			res, err = us.GetByIdProto(context.Background(), &gen.GetByIdReq{Id: tc.Id})
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, res.User.Id, tc.Id)

		})
	}

}

func TestRPCOneWay(t *testing.T) {
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
	var c *Client
	c, err = NewClient("tcp", "127.0.0.1:8081", time.Minute)
	if err != nil {
		t.Log(err)
	}
	err = c.InitService(us)
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
			name: "one way",
			mock: func() {
				u.Msg = "hello world"
				u.Err = errors.New("error")
			},
			wantErr: errors.New("这是单向调用，没有返回值"),
			wantResp: &user_service.GetByIdResponse{
				Data: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			ctx := OneWayContext(context.Background())
			var res *user_service.GetByIdResponse
			res, err = us.GetById(ctx, &user_service.GetByIdRequest{Id: 123})
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, res, tc.wantResp)
			time.Sleep(time.Minute)
		})
	}

}

func TestRPCTimeout(t *testing.T) {
	s, err := NewServer("tcp", "127.0.0.1:8081")
	if err != nil {
		t.Log(err)
	}

	u := &user_service.UserServiceServerTimeout{}
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
	var c *Client
	c, err = NewClient("tcp", "127.0.0.1:8081", time.Minute)
	if err != nil {
		t.Log(err)
	}
	err = c.InitService(us)
	if err != nil {
		t.Log(err)
	}

	testCases := []struct {
		name     string
		mock     func() context.Context
		wantResp *user_service.GetByIdResponse
		wantErr  error
	}{

		{
			name: "timeout",
			mock: func() context.Context {
				u.Id = 123
				u.Err = errors.New("error")
				// 服务睡眠两秒；但是超时设置了一秒
				u.Sleep = time.Second * 2
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				return ctx
			},
			wantErr: context.DeadlineExceeded,
			wantResp: &user_service.GetByIdResponse{
				Data: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.mock()
			var res *user_service.GetByIdResponse
			res, err = us.GetById(ctx, &user_service.GetByIdRequest{Id: 123})
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, res, tc.wantResp)
		})
	}

}

type mockCompressor struct {
}

func (m mockCompressor) Code() uint8 {
	return 2
}

func (m mockCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func (m mockCompressor) UnCompress(data []byte) ([]byte, error) {
	return data, nil
}

func TestRPCCompressor(t *testing.T) {
	s, err := NewServer("tcp", "127.0.0.1:8081")
	if err != nil {
		t.Log(err)
	}

	u := &user_service.UserServiceServerCompressor{}
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

	testCases := []struct {
		name    string
		Id      int
		mock    func()
		wantErr error
		wantRes *user_service.GetByIdResponse
	}{
		{
			name: "normal",
			Id:   123,
			mock: func() {
				var c *Client

				c, err = NewClient("tcp", "127.0.0.1:8081", time.Minute)

				if err != nil {
					t.Log(err)
				}
				err = c.InitService(us)
				if err != nil {
					t.Log(err)
				}
			},
			wantRes: &user_service.GetByIdResponse{
				Data: "123",
			},
		},
		{
			name:    "客户端注册,服务端未注册compressor",
			Id:      123,
			wantErr: errors.New("客户端指定的压缩算法服务端不存在"),
			wantRes: &user_service.GetByIdResponse{Data: ""},
			mock: func() {
				var c *Client

				c, err = NewClient("tcp", "127.0.0.1:8081", time.Minute, ClientWithCompressor(&mockCompressor{}))

				if err != nil {
					t.Log(err)
				}
				err = c.InitService(us)
				if err != nil {
					t.Log(err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			var res *user_service.GetByIdResponse
			res, err = us.GetById(context.Background(), &user_service.GetByIdRequest{Id: tc.Id})
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, res, tc.wantRes)
		})
	}

}
