package round_robin

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"testing"
)

func TestBalancer_Pick(t *testing.T) {
	testCases := []struct {
		name            string
		b               *Balancer
		wantErr         error
		wantBalancerIdx int32
		wantSubConn     *SubConn
	}{
		{
			name: "start",
			b: &Balancer{
				index: -1,
				connections: []balancer.SubConn{
					SubConn{Name: "127.0.0.1:8080"},
					SubConn{Name: "127.0.0.1:8081"},
				},
				len: 2,
			},
			wantBalancerIdx: 0,
			wantSubConn: &SubConn{
				Name: "127.0.0.1:8080",
			},
		},
		{
			name: "next",
			b: &Balancer{
				index: 0,
				connections: []balancer.SubConn{
					SubConn{Name: "127.0.0.1:8080"},
					SubConn{Name: "127.0.0.1:8081"},
				},
				len: 2,
			},
			wantBalancerIdx: 1,
			wantSubConn: &SubConn{
				Name: "127.0.0.1:8081",
			},
		},
		{
			name: "end go start",
			b: &Balancer{
				index: 1,
				connections: []balancer.SubConn{
					SubConn{Name: "127.0.0.1:8080"},
					SubConn{Name: "127.0.0.1:8081"},
				},
				len: 2,
			},
			wantBalancerIdx: 2,
			wantSubConn: &SubConn{
				Name: "127.0.0.1:8080",
			},
		},
		{
			name: "no connection",
			b: &Balancer{
				index:       0,
				connections: []balancer.SubConn{},
				len:         0,
			},
			wantErr: balancer.ErrNoSubConnAvailable,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			balancerRes, err := tc.b.Pick(balancer.PickInfo{})
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, balancerRes.SubConn.(SubConn).Name, tc.wantSubConn.Name)
			assert.NotNil(t, balancerRes.Done)
			assert.Equal(t, tc.wantBalancerIdx, tc.b.index)
		})
	}
}

type SubConn struct {
	Name string
	balancer.SubConn
}
