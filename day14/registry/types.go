package registry

import (
	"context"
	"io"
)

type Registry interface {
	Registry(ctx context.Context, si ServiceInstance) error
	UnRegistry(ctx context.Context, si ServiceInstance) error
	ListServices(ctx context.Context, serviceName string) ([]ServiceInstance, error)
	Subscribe(serviceName string) <-chan Event
	io.Closer
}

type ServiceInstance struct {
	Name string
	//Address resolver.Address // 不能用，里面的attributes中的m拿不到
	Addr   string
	Weight uint32
	Group  string
}

type Event struct {
}
