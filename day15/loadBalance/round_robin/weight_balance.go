package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"sync"
)

// 加权轮询

type WeightBalancer struct {
	connections []*weightConn
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var totalWeight uint32
	var res *weightConn

	for _, connection := range w.connections {
		connection.mutex.Lock()
		totalWeight = totalWeight + connection.efficientWeight
		connection.currentWeight = connection.currentWeight + connection.efficientWeight
		if res == nil || res.currentWeight < connection.currentWeight {
			res = connection
		}
		connection.mutex.Unlock()
	}

	res.mutex.Lock()
	res.currentWeight = res.currentWeight - totalWeight
	res.mutex.Unlock()
	return balancer.PickResult{
		SubConn: res.c,
		Done: func(info balancer.DoneInfo) {
			res.mutex.Lock()
			if info.Err != nil {
				// 不减了，再减溢出边界了
				if res.efficientWeight == uint32(0) {
					return
				}
				res.efficientWeight--
			} else {
				// 不加了，再加溢出边界变成负数或者0了
				if res.efficientWeight == math.MaxUint32 {
					return
				}
				res.efficientWeight++
			}
			res.mutex.Unlock()
		},
	}, nil
}

type WeightBalancerBuilder struct {
}

func (w *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]*weightConn, 0, len(info.ReadySCs))
	for c, ci := range info.ReadySCs {
		weight := ci.Address.Attributes.Value("weight").(uint32)
		connection := &weightConn{
			weight:          weight,
			currentWeight:   weight,
			efficientWeight: weight,
			c:               c,
		}

		connections = append(connections, connection)
	}
	return &WeightBalancer{
		connections: connections,
	}
}

type weightConn struct {
	c               balancer.SubConn
	weight          uint32
	currentWeight   uint32
	efficientWeight uint32
	mutex           sync.Mutex
}
