package user_service

import (
	"context"
	"geektime-go/day13/serialize/proto/gen"
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
