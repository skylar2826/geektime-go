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
	c       *clientv3.Client
	sess    *concurrency.Session
	cancels []func()
	mutex   sync.Mutex
}

func NewRegistry(c *clientv3.Client) (*Registry, error) {
	sess, err := concurrency.NewSession(c) // etcd 提供的租约
	if err != nil {
		return nil, err
	}
	return &Registry{
		c:    c,
		sess: sess,
	}, nil
}

func (r *Registry) Registry(ctx context.Context, instance registry.ServiceInstance) error {
	val, err := json.Marshal(instance)
	if err != nil {
		return err
	}
	instanceKey := r.instanceKey(instance)
	_, err = r.c.Put(ctx, instanceKey, string(val), clientv3.WithLease(r.sess.Lease()))
	return err
}

func (r *Registry) UnRegistry(ctx context.Context, instance registry.ServiceInstance) error {
	_, err := r.c.Delete(ctx, r.instanceKey(instance))
	return err
}

func (r *Registry) ListServices(ctx context.Context, serviceName string) ([]registry.ServiceInstance, error) {
	getResp, err := r.c.Get(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	res := make([]registry.ServiceInstance, 0, len(getResp.Kvs))
	for _, kv := range getResp.Kvs {
		var si registry.ServiceInstance
		err = json.Unmarshal(kv.Value, &si)
		if err != nil {
			return nil, err
		}
		res = append(res, si)
	}
	return res, nil
}

func (r *Registry) Subscribe(serviceName string) (<-chan registry.Event, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r.mutex.Lock()
	r.cancels = append(r.cancels, cancel)
	r.mutex.Unlock()
	ctx = clientv3.WithRequireLeader(ctx)
	watchResp := r.c.Watch(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
	res := make(chan registry.Event)

	go func() {
		for {
			select {
			case resp := <-watchResp:
				if resp.Err() != nil || resp.Canceled {
					return
				}
				for range resp.Events {
					res <- registry.Event{}
				}
			case <-ctx.Done():
				return
			}

		}
	}()
	return res, nil
}

func (r *Registry) Close() error {
	r.mutex.Lock()
	cancels := r.cancels
	r.cancels = nil
	r.mutex.Unlock()
	for _, cancel := range cancels {
		cancel()
	}
	return r.sess.Close()
}

func (r *Registry) instanceKey(instance registry.ServiceInstance) string {
	return fmt.Sprintf("/micro/%s/%s", instance.Name, instance.Address)
}

func (r *Registry) serviceKey(instanceName string) string {
	return fmt.Sprintf("/micro/%s", instanceName)
}
