package user_service

import (
	"context"
	"geektime-go/day13/serialize/proto/gen"
)

type UserService struct {
	GetById      func(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error)
	GetByIdProto func(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error)
}

func (u *UserService) Name() string {
	return "user_service"
}
