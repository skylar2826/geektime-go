package user_service

import (
	"context"
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
