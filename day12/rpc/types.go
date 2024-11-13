package rpc

import "context"

type Request struct {
	ServiceName string `json:"serviceName"`
	MethodName  string `json:"methodName"`
	Arg         []byte `json:"arg"`
}

type Response struct {
	Data []byte `json:"data"`
}

type proxy interface {
	invoke(ctx context.Context, request *Request) (*Response, error)
}

type Service interface {
	Name() string
}
