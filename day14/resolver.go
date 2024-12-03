package day14

import (
	"context"
	"geektime-go/day14/registry/etcd"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"time"
)

type ResolverBuilder struct {
	registry *etcd.Registry
	timeout  time.Duration
}

type BuilderOption func(r *ResolverBuilder)

func BuilderWithTimeout(timeout time.Duration) BuilderOption {
	return func(r *ResolverBuilder) {
		r.timeout = timeout
	}
}

func NewResolverBuilder(registry *etcd.Registry, opts ...BuilderOption) *ResolverBuilder {
	r := &ResolverBuilder{
		registry: registry,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	resolverInstance := &Resolver{
		registry: r.registry,
		target:   target,
		cc:       cc,
		timeout:  r.timeout,
	}
	resolverInstance.Resolve()
	go resolverInstance.watch()
	return resolverInstance, nil
}

func (r *ResolverBuilder) Scheme() string {
	return "passthrough"
}

type Resolver struct {
	registry *etcd.Registry
	timeout  time.Duration
	target   resolver.Target
	cc       resolver.ClientConn
	close    chan struct{}
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	r.Resolve()
}

func (r *Resolver) watch() {
	eventChan := r.registry.Subscribe(r.target.Endpoint())
	for {
		select {
		case <-eventChan:
			r.Resolve()
		case <-r.close:
			return
		}
	}
}

func (r *Resolver) Resolve() {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	serviceInstances, err := r.registry.ListServices(ctx, r.target.Endpoint())
	if err != nil {
		r.cc.ReportError(err)
		return
	}
	addresses := make([]resolver.Address, 0, len(serviceInstances))
	for _, si := range serviceInstances {
		addresses = append(addresses, resolver.Address{
			Addr:       si.Addr,
			Attributes: attributes.New("Weight", si.Weight).WithValue("Group", si.Group),
		})
	}
	err = r.cc.UpdateState(resolver.State{
		Addresses: addresses,
	})
	if err != nil {
		r.cc.ReportError(err)
		return
	}
}

func (r *Resolver) Close() {
	close(r.close)
}
