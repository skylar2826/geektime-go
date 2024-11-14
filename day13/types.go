package rpc

import (
	"context"
	"geektime-go/day13/message"
)

type proxy interface {
	invoke(ctx context.Context, request *message.Request) (*message.Response, error)
}

type Service interface {
	Name() string
}
