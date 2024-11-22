package registry

import (
	"context"
	"io"
)

type Registry interface {
	Registry(ctx context.Context, instance ServiceInstance) error
	UnRegistry(ctx context.Context, instance ServiceInstance) error
	ListServices(ctx context.Context, serviceName string) ([]ServiceInstance, error)
	Subscribe(serviceName string) (<-chan Event, error)
	io.Closer
}

type ServiceInstance struct {
	Address string
	Name    string
}

type Event struct {
}
