package random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand/v2"
)

type WeightBalancer struct {
	totalWeight uint32
	connections []weightConn
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	randWeight := rand.IntN(int(w.totalWeight + 1)) // [0, n]
	for _, c := range w.connections {
		if int(c.weight) >= randWeight {
			return balancer.PickResult{SubConn: c.c}, nil
		}
	}
	// 一定会匹配，不会走这里
	panic("panic")
}

type WeightBalancerBuilder struct {
}

func (w *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]weightConn, 0, len(info.ReadySCs))
	var totalWeight uint32
	for c, ci := range info.ReadySCs {
		weight := ci.Address.Attributes.Value("Weight").(uint32)
		totalWeight += weight
		connections = append(connections, weightConn{
			c:      c,
			weight: weight,
		})
	}
	return &WeightBalancer{
		connections: connections,
		totalWeight: totalWeight,
	}
}

type weightConn struct {
	c      balancer.SubConn
	weight uint32
}
