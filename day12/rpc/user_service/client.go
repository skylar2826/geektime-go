package user_service

import (
	"context"
)

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error)
}

func (u *UserService) Name() string {
	return "user_service"
}
