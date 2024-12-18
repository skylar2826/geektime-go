package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync/atomic"
)

// picker 是纯粹的负载均衡算法
// balancer 是对picker的封装，里面有与clientConn打交道

// 轮询

type Balancer struct {
	index       int32
	connections []balancer.SubConn
	len         int32
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	idx := atomic.AddInt32(&b.index, 1)
	c := b.connections[idx%b.len]

	return balancer.PickResult{
		// SubConn 对一个实例的连接池的抽象
		// clientConn 是对一个服务连接池的抽象，clientConn与SubConn是一对多的关系
		SubConn: c,
		// 响应回来的回调
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type Builder struct {
}

func (b Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for c := range info.ReadySCs {
		connections = append(connections, c)
	}
	return &Balancer{connections: connections, index: -1, len: int32(len(info.ReadySCs))}
}
