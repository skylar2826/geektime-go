package round_robin

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

func TestE2EBalancer_Pick(t *testing.T) {
	go func() {
		us := &userServiceServer{}
		server := grpc.NewServer()
		gen.RegisterUserServiceServer(server, us)
		l, err := net.Listen("tcp", "127.0.0.1:8085")
		require.NoError(t, err)
		err = server.Serve(l)
		t.Log(err)
	}()

	time.Sleep(time.Second * 3)
	builder := base.NewBalancerBuilder("TEST_DEMO_ROUND_ROBIN", &Builder{}, base.Config{HealthCheck: true})
	balancer.Register(builder)
	conn, err := grpc.Dial("127.0.0.1:8085", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"TEST_DEMO_ROUND_ROBIN"}`))
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
