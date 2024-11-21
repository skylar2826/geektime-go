package user_service

import (
	"context"
	"fmt"
	"geektime-go/day13/serialize/proto/gen"
	"testing"
	"time"
)

type UserServiceServer struct {
	Msg string
	Err error
}

func (u *UserServiceServer) Name() string {
	return "user_service"
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error) {
	return &GetByIdResponse{
		Data: u.Msg,
	}, u.Err
}

func (u *UserServiceServer) GetByIdProto(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	return &gen.GetByIdResp{
		User: &gen.User{
			Id: req.Id,
		},
	}, nil
}

type UserServiceServerTimeout struct {
	t     testing.T
	Sleep time.Duration
	Err   error
	Id    int
}

func (u *UserServiceServerTimeout) Name() string {
	return "user_service"
}

func (u *UserServiceServerTimeout) GetById(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error) {
	if _, ok := ctx.Deadline(); !ok {
		u.t.Fatal("没有设置超时时间")
	}
	time.Sleep(u.Sleep)
	return &GetByIdResponse{
		Data: fmt.Sprint(req.Id),
	}, nil
}

type UserServiceServerCompressor struct {
}

func (u *UserServiceServerCompressor) Name() string {
	return "user_service"
}

func (u *UserServiceServerCompressor) GetById(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error) {
	return &GetByIdResponse{
		Data: fmt.Sprint(req.Id),
	}, nil
}
