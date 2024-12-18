package round_robin

import (
	"geektime-go/day15/route"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"sync/atomic"
)

// picker 是纯粹的负载均衡算法
// balancer 是对picker的封装，里面有与clientConn打交道

// 轮询

type Balancer struct {
	index       int32
	connections []subConn
	len         int32
	filter      route.Filter
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	candidates := make([]subConn, 0, len(b.connections))
	for _, c := range b.connections {
		if b.filter != nil && !b.filter(info, c.addr) {
			continue
		}
		candidates = append(candidates, c)
	}

	if len(candidates) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	idx := atomic.AddInt32(&b.index, 1)
	c := candidates[int(idx)%len(candidates)]

	return balancer.PickResult{
		// SubConn 对一个实例的连接池的抽象
		// clientConn 是对一个服务连接池的抽象，clientConn与SubConn是一对多的关系
		SubConn: c.c,
		// 响应回来的回调
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type Builder struct {
	Filter route.Filter
}

func (b Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]subConn, 0, len(info.ReadySCs))
	for c, ci := range info.ReadySCs {
		connections = append(connections, subConn{
			c:    c,
			addr: ci.Address,
		})
	}

	return &Balancer{connections: connections, index: -1, len: int32(len(info.ReadySCs)), filter: b.Filter}
}

type subConn struct {
	c    balancer.SubConn
	addr resolver.Address
}
