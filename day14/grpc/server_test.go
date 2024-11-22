package grpc

import (
	"context"
	"fmt"
	"geektime-go/day13/serialize/proto/gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	us := &Server{}
	server := grpc.NewServer()                // 注册中心
	gen.RegisterUserServiceServer(server, us) // 注册服务
	listen, err := net.Listen("tcp", "127.0.0.1:8081")
	require.NoError(t, err)
	err = server.Serve(listen) // 启动
	t.Log(err)
}

type Server struct {
	gen.UnimplementedUserServiceServer
}

func (s *Server) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	fmt.Println(req)
	return &gen.GetByIdResp{
		User: &gen.User{Id: 123, Name: "hello world"},
	}, nil
}
