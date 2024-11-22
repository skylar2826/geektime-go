package day14

import (
	"context"
	"geektime-go/day13/serialize/proto/gen"

	"geektime-go/day14/registry/etcd"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
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

	c := NewClient(ClientInSecure(), ClientWithRegistry(registry, time.Second*3))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clientConn, err := c.Dial(ctx, "user-service")
	require.NoError(t, err)

	uc := gen.NewUserServiceClient(clientConn)

	var resp *gen.GetByIdResp
	resp, err = uc.GetById(ctx, &gen.GetByIdReq{Id: 13})
	require.NoError(t, err)
	t.Log(resp)
}
