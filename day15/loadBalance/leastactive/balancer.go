package leastactive

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"sync/atomic"
)

// 请求活跃数

type Balancer struct {
	connections []activeConn
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	res := activeConn{
		cnt: math.MaxUint32,
	}

	for _, c := range b.connections {
		// 找到第一个请求活跃数小于res的
		if atomic.LoadUint32(&c.cnt) <= res.cnt {
			res = c
		}
	}
	atomic.AddUint32(&res.cnt, 1)
	return balancer.PickResult{
		SubConn: res.c,
		Done: func(info balancer.DoneInfo) {
			atomic.AddUint32(&res.cnt, -1)
		},
	}, nil
}

type Builder struct {
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]activeConn, 0, len(info.ReadySCs))
	for c := range info.ReadySCs {
		connections = append(connections, activeConn{
			c: c,
		})
	}
	return &Balancer{
		connections: connections,
	}
}

type activeConn struct {
	cnt uint32
	c   balancer.SubConn
}
