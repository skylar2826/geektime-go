package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"geektime-go/day14/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync"
)

type Registry struct {
	cc      *clientv3.Client
	sess    *concurrency.Session
	mutex   sync.Mutex
	cancels []func()
}

func NewRegistry(cc *clientv3.Client) (*Registry, error) {
	sess, err := concurrency.NewSession(cc)
	if err != nil {
		return nil, err
	}
	r := &Registry{
		cc:   cc,
		sess: sess,
	}
	return r, nil
}

func (r *Registry) Registry(ctx context.Context, si registry.ServiceInstance) error {
	val, err := json.Marshal(si)
	if err != nil {
		return err
	}

	fmt.Println("instanceKey", r.instanceKey(si))
	_, err = r.cc.Put(ctx, r.instanceKey(si), string(val), clientv3.WithLease(r.sess.Lease()))
	return err
}

func (r *Registry) UnRegistry(ctx context.Context, si registry.ServiceInstance) error {
	_, err := r.cc.Delete(ctx, r.instanceKey(si))
	return err
}

func (r *Registry) ListServices(ctx context.Context, serviceName string) ([]registry.ServiceInstance, error) {
	resp, err := r.cc.Get(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	si := make([]registry.ServiceInstance, len(resp.Kvs))
	for idx, kv := range resp.Kvs {
		err = json.Unmarshal(kv.Value, &si[idx])
		if err != nil {
			return nil, err
		}
	}
	return si, nil
}

func (r *Registry) Subscribe(serviceName string) <-chan registry.Event {
	ctx, cancel := context.WithCancel(context.Background())
	r.mutex.Lock()
	r.cancels = append(r.cancels, cancel)
	r.mutex.Unlock()
	watchChan := r.cc.Watch(ctx, r.serviceKey(serviceName))
	res := make(chan registry.Event)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case resp := <-watchChan:
				if resp.Canceled || resp.Err() != nil {
					return
				}
				for range resp.Events {
					res <- registry.Event{}
				}
			}
		}
	}()
	return res
}

func (r *Registry) Close() error {
	cancels := r.cancels
	r.mutex.Lock()
	r.cancels = nil
	r.mutex.Unlock()
	for _, cancel := range cancels {
		cancel()
	}
	return r.sess.Close()
}

func (r *Registry) instanceKey(si registry.ServiceInstance) string {
	return fmt.Sprintf("/micro/%s/%s", si.Name, si.Addr)
}

func (r *Registry) serviceKey(serviceName string) string {
	return fmt.Sprintf("/micro/%s", serviceName)
}
