package broadcast

import (
	"context"
	"fmt"
	day14 "geektime-go/day14"
	"geektime-go/day14/proto/gen"
	"geektime-go/day14/registry/etcd"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
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

	var servers []*userServiceServer
	var eg errgroup.Group
	for i := 0; i < 5; i++ {
		s := day14.NewServer("user-service", fmt.Sprintf("127.0.0.1:808%d", i), day14.ServerWithRegistry(registry, time.Second*10))
		us := &userServiceServer{
			idx: i,
		}
		servers = append(servers, us)
		gen.RegisterUserServiceServer(s, us)
		eg.Go(func() error {
			return s.Start()
		})
		defer func() {
			_ = s.Close()
		}()
	}

	if err != nil {
		t.Log(err)
	}

	time.Sleep(time.Second * 3)

	c := day14.NewClient(day14.ClientWithInsecure(), day14.ClientWithRegistry(registry, time.Second*3000))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()
	var conn *grpc.ClientConn
	var respChan <-chan Resp
	ctx, respChan = useBroadcast(ctx)
	//接收所有响应
	//go func() {
	//	for res := range respChan {
	//		fmt.Println(res.Reply, res.Err)
	//	}
	//}()
	// 接收第一个响应
	go func() {
		res := <-respChan
		fmt.Println(res.Reply, res.Err)
	}()
	broadcast := NewClusterBuilder(registry, "user-service", grpc.WithInsecure())
	conn, err = c.Dial(ctx, "user-service", grpc.WithUnaryInterceptor(broadcast.BuildUnaryInterceptor()))
	if err != nil {
		t.Log(err)
	}
	uc := gen.NewUserServiceClient(conn)
	var resp *gen.GetByIdResp

	resp, err = uc.GetById(ctx, &gen.GetByIdReq{Id: 123})
	if err != nil {
		t.Log(err)
	}
	t.Log(resp)

	for _, s := range servers {
		fmt.Println("服务调用次数", s.cnt)
		assert.Equal(t, s.cnt, 1)
	}
}

type userServiceServer struct {
	gen.UnimplementedUserServiceServer
	cnt int
	idx int
}

func (u *userServiceServer) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	u.cnt++

	return &gen.GetByIdResp{User: &gen.User{Id: req.Id, Name: fmt.Sprintf("hello world %d", u.idx)}}, nil
}
