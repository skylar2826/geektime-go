package day14

import (
	"context"
	"fmt"
	"geektime-go/day14/proto/gen"
	"geektime-go/day14/registry/etcd"
	"geektime-go/day15/route"
	"geektime-go/day15/route/round_robin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	cc, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	var registry *etcd.Registry
	registry, err = etcd.NewRegistry(cc)
	if err != nil {
		t.Log(err)
	}
	c := NewClient(ClientWithInsecure(), ClientWithRegistry(registry, time.Second*3000), ClientWithPickBuilder("GROUP_ROUND_ROBIN", &round_robin.Builder{
		Filter: (route.GroupFilterBuilder{
			Group: "B",
		}).Build(),
	}))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()
	var conn *grpc.ClientConn
	ctx = context.WithValue(ctx, "Group", "B")
	conn, err = c.Dial(ctx, "user-service")
	if err != nil {
		t.Log(err)
	}
	uc := gen.NewUserServiceClient(conn)
	var resp *gen.GetByIdResp
	for i := 0; i < 10; i++ {
		fmt.Println("准备发送请求")
		resp, err = uc.GetById(ctx, &gen.GetByIdReq{Id: 123})
		if err != nil {
			t.Log(err)
		}
		t.Log(resp)
	}

}
