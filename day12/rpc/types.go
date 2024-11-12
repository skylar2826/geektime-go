package rpc

import "context"

type Service interface {
	Name() string
}

type Request struct {
	ServiceName string
	MethodName  string
	Arg         []byte
}

type Proxy interface {
	invoke(ctx context.Context, req *Request) (*Response, error)
}

type Response struct {
	data []byte
}
