package round_robig

import (
	"context"
	"fmt"
	"geektime-go/day14/proto/gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"net"
	"testing"
	"time"
)

func TestBalancer(t *testing.T) {
	go func() {
		us := &userServiceServer{}
		server := grpc.NewServer()
		gen.RegisterUserServiceServer(server, us)
		l, err := net.Listen("tcp", ":8081")
		require.NoError(t, err)
		err = server.Serve(l)
		t.Log(err)
	}()

	time.Sleep(time.Second * 3)
	balancer.Register(base.NewBalancerBuilder("DEMO_ROUND_ROBIG", &Builder{}, base.Config{HealthCheck: true}))
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure(), grpc.WithDefaultServiceConfig())
	require.NoError(t, err)
	uc := gen.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var resp *gen.GetByIdResp
	resp, err = uc.GetById(ctx, &gen.GetByIdReq{
		Id: 123,
	})
	require.NoError(t, err)
	t.Log(resp)
}

type userServiceServer struct {
	gen.UnimplementedUserServiceServer
}

func (u userServiceServer) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	fmt.Println(req)
	return &gen.GetByIdResp{User: &gen.User{Id: req.Id, Name: "hello world"}}, nil
}
