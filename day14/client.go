package day14

import (
	"context"
	"encoding/json"
	"fmt"
	"geektime-go/day14/registry/etcd"
	"geektime-go/day15/route/round_robin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"time"
)

type Client struct {
	insecure        bool
	registry        *etcd.Registry
	timeout         time.Duration
	pickBuilderName string
}

type ClientOption func(c *Client)

func ClientWithInsecure() ClientOption {
	return func(c *Client) {
		c.insecure = true
	}
}

func ClientWithRegistry(registry *etcd.Registry, timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.registry = registry
		c.timeout = timeout
	}
}

func ClientWithPickBuilder(name string, pb *round_robin.Builder) ClientOption {
	balanceBuilder := base.NewBalancerBuilder(name, pb, base.Config{HealthCheck: true})
	balancer.Register(balanceBuilder)

	return func(c *Client) {
		c.pickBuilderName = name
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) dial(ctx context.Context, serviceName string) (*grpc.ClientConn, error) {
	dialOptions := make([]grpc.DialOption, 0, 4)
	if c.insecure {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}
	if c.registry != nil {
		resolverBuilder := NewResolverBuilder(c.registry, BuilderWithTimeout(c.timeout))
		dialOptions = append(dialOptions, grpc.WithResolvers(resolverBuilder))
	}
	if c.pickBuilderName != "" {
		bs, err := json.Marshal(map[string]string{"loadBalancingPolicy": c.pickBuilderName})
		if err != nil {
			return nil, err
		}
		dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(string(bs)))
	}
	//conn, err := grpc.Dial(c.address, dialOptions...)
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("passthrough:///%s", serviceName), dialOptions...)
	return conn, err
}
