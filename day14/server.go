package day14

import (
	"context"
	"geektime-go/day14/registry"
	"geektime-go/day14/registry/etcd"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	name    string
	address string
	*grpc.Server
	registry *etcd.Registry
	timeout  time.Duration
	weight   uint32
	group    string
}

type ServerOption func(s *Server)

func ServerWithRegistry(registry *etcd.Registry, timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.registry = registry
		s.timeout = timeout
	}
}

// ServerWithWeight 服务权重
func ServerWithWeight(weight uint32) ServerOption {
	return func(s *Server) {
		s.weight = weight
	}
}

func ServerWithGroup(group string) ServerOption {
	return func(s *Server) {
		s.group = group
	}
}

func NewServer(name string, address string, opts ...ServerOption) *Server {
	s := &Server{name: name, address: address, Server: grpc.NewServer()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)

	if err != nil {
		return err
	}
	if s.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		defer cancel()
		err = s.registry.Registry(ctx, registry.ServiceInstance{
			Name:   s.name,
			Addr:   listener.Addr().String(),
			Weight: s.weight,
			Group:  s.group,
		})
		if err != nil {
			return err
		}
	}
	err = s.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	if s.registry != nil {
		err := s.registry.Close()
		if err != nil {
			return err
		}
	}
	s.GracefulStop()
	return nil
}
