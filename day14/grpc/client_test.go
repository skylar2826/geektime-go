package grpc

import (
	"context"
	"geektime-go/day13/serialize/proto/gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	clientConn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	require.NoError(t, err)
	client := gen.NewUserServiceClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var resp *gen.GetByIdResp
	resp, err = client.GetById(ctx, &gen.GetByIdReq{Id: 13})
	require.NoError(t, err)
	t.Log(resp)
}
