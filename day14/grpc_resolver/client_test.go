package grpc

import (
	"context"
	"geektime-go/day13/serialize/proto/gen"
	"geektime-go/day14/registry/etcd"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	registry, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)

	registryBuilder, err := NewRegistryBuilder(registry)
	require.NoError(t, err)

	clientConn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure(), grpc.WithResolvers(registryBuilder))



	require.NoError(t, err)
	client := gen.NewUserServiceClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var resp *gen.GetByIdResp
	resp, err = client.GetById(ctx, &gen.GetByIdReq{Id: 13})
	require.NoError(t, err)
	t.Log(resp)
}
