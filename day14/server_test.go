package day14

import (
	"context"
	"fmt"
	"geektime-go/day14/proto/gen"
	"geektime-go/day14/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	cc, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		t.Log(err)
	}
	var registry *etcd.Registry
	registry, err = etcd.NewRegistry(cc)
	if err != nil {
		t.Log(err)
	}

	var eg errgroup.Group
	for i := 0; i < 5; i++ {
		var group = "A"
		if i%2 == 0 {
			group = "B"
		}
		s := NewServer("user-service", fmt.Sprintf("127.0.0.1:808%d", i), ServerWithRegistry(registry, time.Second*10), ServerWithGroup(group), ServerWithWeight(uint32((i+1)*5)))
		us := &userServiceServer{
			group: group,
			port:  fmt.Sprintf("808%d", i),
		}
		gen.RegisterUserServiceServer(s, us)
		eg.Go(func() error {
			return s.Start()
		})
	}
	err = eg.Wait()
	if err != nil {
		t.Log(err)
	}
}

type userServiceServer struct {
	gen.UnimplementedUserServiceServer
	group string
	port  string
}

func (u userServiceServer) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	//fmt.Println(req)
	fmt.Println(u.group, u.port)

	return &gen.GetByIdResp{User: &gen.User{Id: req.Id, Name: "hello world"}}, nil
}
