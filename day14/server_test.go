package day14

import (
	"context"
	"fmt"
	"geektime-go/day13/serialize/proto/gen"

	"geektime-go/day14/registry/etcd"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

func TestServer(t *testing.T) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	}) // get client
	require.NoError(t, err)
	var r *etcd.Registry
	r, err = etcd.NewRegistry(c) // get registry
	require.NoError(t, err)
	server := NewServer("user-service", ServerWithRegistry(r)) // get server
	//server := grpc.NewServer()
	us := &UserServiceServer{}
	gen.RegisterUserServiceServer(server, us)
	err = server.Start("127.0.0.1:8081") // 启动(含：注册服务)
	t.Log(err)
}

type UserServiceServer struct {
	gen.UnimplementedUserServiceServer
}

func (s *UserServiceServer) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	fmt.Println(req)
	return &gen.GetByIdResp{
		User: &gen.User{Id: 123, Name: "hello world"},
	}, nil
}
