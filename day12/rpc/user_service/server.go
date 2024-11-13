package user_service

import (
	"context"
	"fmt"
)

type UserServiceServer struct {
}

func (u *UserServiceServer) Name() string {
	return "user_service"
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error) {
	str := fmt.Sprintf("请求Id: %d, 响应信息：hello world", req.Id)
	return &GetByIdResponse{
		Data: str,
	}, nil
}
