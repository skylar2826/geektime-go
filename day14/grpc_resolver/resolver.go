package grpc

import (
	"context"
	"geektime-go/day14/registry"
	"google.golang.org/grpc/resolver"
	"time"
)

type GrpcResolverBuilder struct {
	registry registry.Registry
	timeout  time.Duration
}

type BuilderOption func(b *GrpcResolverBuilder)

func BuilderWithTimeout(timeout time.Duration) BuilderOption {
	return func(b *GrpcResolverBuilder) {
		b.timeout = timeout
	}
}

func NewRegistryBuilder(r registry.Registry, opts ...BuilderOption) (*GrpcResolverBuilder, error) {
	g := &GrpcResolverBuilder{
		registry: r,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g, nil
}

func (g *GrpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	res := &grpcResolver{
		cc:      cc,
		r:       g.registry,
		target:  target,
		timeout: g.timeout,
	}
	res.Resolve() // 立刻解析一次；不然客户端不知道怎么连接，客户端就永远连接不上；而且等待客户连，容易timeout
	go res.watch()
	return res, nil
}

func (g *GrpcResolverBuilder) Scheme() string {
	return "passthrough"
}

type grpcResolver struct {
	cc      resolver.ClientConn
	r       registry.Registry
	target  resolver.Target
	timeout time.Duration
	close   chan struct{}
}

func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	g.Resolve()
}

func (g *grpcResolver) watch() {
	events, err := g.r.Subscribe(g.target.Endpoint())
	if err != nil {
		g.cc.ReportError(err)
		return
	}
	for {
		select {
		case <-events:
			g.Resolve()
		case <-g.close:
			return
		}
	}
}

func (g *grpcResolver) Resolve() {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	instances, err := g.r.ListServices(ctx, g.target.Endpoint())
	if err != nil {
		g.cc.ReportError(err)
		return
	}
	address := make([]resolver.Address, 0, len(instances))
	for _, si := range instances {
		address = append(address, resolver.Address{Addr: si.Address})
	}
	err = g.cc.UpdateState(resolver.State{
		Addresses: address,
	})
	if err != nil {
		g.cc.ReportError(err)
		return
	}
}

func (g *grpcResolver) Close() {
	close(g.close)
}
