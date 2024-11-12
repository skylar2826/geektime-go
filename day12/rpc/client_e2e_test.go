package rpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInitClientProxy(t *testing.T) {
	go func() {
		server := NewServer()
		server.registerService(&UserServiceServer{})
		err := server.Start("tcp", "127.0.0.1:8081")
		t.Log(err)
	}()

	time.Sleep(time.Second * 5)
	usClient := &UserService{}
	err := InitClientProxy("127.0.0.1:8081", usClient)
	t.Log(err)
	var resp *GetByIdResp
	resp, err = usClient.GetById(context.Background(), &GetByIdReq{Id: "123"})
	t.Log(err)
	assert.Equal(t, resp, &GetByIdResp{
		str: "hello world"},
	)
}
