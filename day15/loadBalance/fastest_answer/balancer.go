package fastest_answer

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// 快响应负载均衡算法 响应时间

type Balancer struct {
	connections []*conn
	mutex       sync.RWMutex
	lastSync    time.Time
	endpoint    string
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	b.mutex.RLocker()
	if len(b.connections) == 0 {
		b.mutex.RUnlock()
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var res *conn
	for _, c := range b.connections {
		if res == nil || res.responseTime > c.responseTime {
			res = c
		}
	}
	b.mutex.RUnlock()
	return balancer.PickResult{
		SubConn: res.c,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type Builder struct {
	Interval time.Duration
	Endpoint string // prometheus 的地址
	Query    string
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]*conn, 0, len(info.ReadySCs))
	for c, val := range info.ReadySCs {
		connections = append(connections, &conn{
			c:    c,
			addr: val.Address,
			// 随意设置一个默认值。当然这个默认值会对初始的负载均衡有影响
			// 不过一段时间后就没影响了
			responseTime: time.Millisecond * 100,
		})
	}

	res := &Balancer{
		connections: connections,
	}

	// 这里有一个很大的问题：不好退出，因为没有grpc 不会调用 Close 方法
	// 可以考虑用 runtime.SetFinalizer 来在 res 被回收时得到通知
	ch := make(chan struct{}, 1)
	runtime.SetFinalizer(res, func() {
		ch <- struct{}{}
	})

	go func() {
		ticker := time.NewTicker(b.Interval)
		for {
			select {
			case <-ticker.C:
				// 这里很难容错，即如果刷新响应时间失败了咋整
				res.updateRespTime(b.Endpoint, b.Query)
			case <-ch:
				return
			}
		}
	}()

	return res
}

func (b *Balancer) updateRespTime(endpoint, query string) {
	// 这里很难容错，即如果刷新响应时间失败咋整
	httpResp, err := http.Get(fmt.Sprintf("%s/api/v1/query?query=%s", endpoint, query))
	if err != nil {
		log.Fatal("查询 prometheus 失败", err)
		return
	}
	decoder := json.NewDecoder(httpResp.Body)

	var resp response
	err = decoder.Decode(&resp)
	if err != nil {
		log.Fatal("反序列化 http 响应失败", err)
		return
	}

	for _, promRes := range resp.Data.Result {
		addr, ok := promRes.Metric["address"]
		if !ok {
			return
		}

		for _, c := range b.connections {
			if c.addr.Addr == addr {
				ms, er := strconv.ParseInt(promRes.Value[1].(string), 10, 64)
				if er != nil {
					continue
				}
				c.responseTime = time.Duration(ms) * time.Millisecond
			}
		}
	}
}

type conn struct {
	c            balancer.SubConn
	responseTime time.Duration
	addr         resolver.Address
}

type response struct {
	Status string `json:"status"`
	Data   data   `json:"data"`
}

type data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}
