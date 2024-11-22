package day14

import (
	"context"
	"fmt"
	grpc2 "geektime-go/day14/grpc_resolver"
	"geektime-go/day14/registry"
	"google.golang.org/grpc"
	"time"
)

type ClientOption func(c *Client)

type Client struct {
	insecure bool
	registry registry.Registry
	timeout  time.Duration
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func ClientInSecure() ClientOption {
	return func(c *Client) {
		c.insecure = true
	}
}

func ClientWithRegistry(registry registry.Registry, timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.registry = registry
		c.timeout = timeout
	}
}

func (c *Client) Dial(ctx context.Context, service string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if c.registry != nil {
		rb, err := grpc2.NewRegistryBuilder(c.registry, grpc2.BuilderWithTimeout(c.timeout))
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithResolvers(rb))
	}
	if c.insecure {
		opts = append(opts, grpc.WithInsecure())
	}

	clientConn, err := grpc.DialContext(ctx, fmt.Sprintf("passthrough:///%s", service), opts...)
	return clientConn, err
}
