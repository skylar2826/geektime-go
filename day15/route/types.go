package route

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

type Filter func(info balancer.PickInfo, addr resolver.Address) bool

type GroupFilterBuilder struct {
	Group string
}

func (g GroupFilterBuilder) Build() Filter {
	return func(info balancer.PickInfo, addr resolver.Address) bool {
		target := addr.Attributes.Value("Group") // server 注册
		input := info.Ctx.Value("Group")         // client 提供
		return input == target
	}
}
