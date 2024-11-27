package round_robig

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type Balancer struct {
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	//TODO implement me
	panic("implement me")
}

type Builder struct {
}

func (b Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	//TODO implement me
	panic("implement me")
}
