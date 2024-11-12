package rpc

import (
	"context"
	"log"
)

// 客户端
type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id string
}

type GetByIdResp struct {
	str string
}

// 服务端
type UserServiceServer struct {
}

func (u *UserServiceServer) Name() string {
	return "user-service"
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	log.Println(req)
	return &GetByIdResp{
		str: "hello world",
	}, nil
}
