package random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand/v2"
)

type Balancer struct {
	connections []balancer.SubConn
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	// [0, n)
	idx := rand.IntN(len(b.connections))
	return balancer.PickResult{
		SubConn: b.connections[idx],
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type Builder struct {
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	//TODO implement me
	panic("implement me")
}
