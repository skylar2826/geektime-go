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
	name            string
	registry        *etcd.Registry
	registryTimeout time.Duration
	*grpc.Server
	listener net.Listener
}

type ServerOption func(s *Server)

func ServerWithRegistry(r *etcd.Registry) ServerOption {
	return func(s *Server) {
		s.registry = r
	}
}

func NewServer(name string, opts ...ServerOption) *Server {
	s := &Server{
		name:            name,
		Server:          grpc.NewServer(),
		registryTimeout: time.Second * 10,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	if s.registry != nil {
		// 注册服务
		ctx, cancel := context.WithTimeout(context.Background(), s.registryTimeout)
		defer cancel()
		err = s.registry.Registry(ctx, registry.ServiceInstance{
			Name:    s.name,
			Address: listener.Addr().String(),
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
	// 先关registry再close listener;因为还有请求要处理
	if s.registry != nil {
		err := s.registry.Close()
		if err != nil {
			return err
		}
	}
	s.GracefulStop()
	//err := s.listener.Close()
	return nil
}
