package round_robin

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"testing"
)

func TestWeightBalancer_Pick(t *testing.T) {
	b := &WeightBalancer{
		connections: []*weightConn{
			&weightConn{
				c: SubConn{
					Name: "weight-5",
				},
				weight:          uint32(5),
				currentWeight:   uint32(5),
				efficientWeight: uint32(5),
			},
			&weightConn{
				c: SubConn{
					Name: "weight-4",
				},
				weight:          uint32(4),
				currentWeight:   uint32(4),
				efficientWeight: uint32(4),
			},
			&weightConn{
				c: SubConn{
					Name: "weight-3",
				},
				weight:          uint32(3),
				currentWeight:   uint32(3),
				efficientWeight: uint32(3),
			},
		},
	}

	res, err := b.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(SubConn).Name, "weight-5")

	res, err = b.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(SubConn).Name, "weight-4")

	res, err = b.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(SubConn).Name, "weight-3")

	res, err = b.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(SubConn).Name, "weight-5")

	res, err = b.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(SubConn).Name, "weight-4")
}
