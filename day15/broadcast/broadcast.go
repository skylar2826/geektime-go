package broadcast

import (
	"context"
	"geektime-go/day14/registry"
	"google.golang.org/grpc"
	"reflect"
	"sync"
)

type ClusterBuilder struct {
	register    registry.Registry
	serviceName string
	opts        []grpc.DialOption
}

func NewClusterBuilder(r registry.Registry, serviceName string, opts ...grpc.DialOption) *ClusterBuilder {
	return &ClusterBuilder{
		register:    r,
		serviceName: serviceName,
		opts:        opts,
	}
}

func (c ClusterBuilder) BuildUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		val, ok := isBroadcast(ctx)
		defer func() {
			close(val)
		}()
		if !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		instances, err := c.register.ListServices(ctx, c.serviceName)
		if err != nil {
			return err
		}
		var wg sync.WaitGroup
		typ := reflect.TypeOf(reply).Elem()
		wg.Add(len(instances))
		for _, instance := range instances {

			go func() {
				iCC, er := grpc.Dial(instance.Addr, c.opts...)
				if er != nil {
					val <- Resp{Err: er}
					wg.Done()
					return
				}
				newReply := reflect.New(typ).Interface()
				er = invoker(ctx, method, req, newReply, iCC, opts...)
				//val <- Resp{
				//	Reply: newReply,
				//	Err:   er,
				//}

				// 所有响应
				// 如果没人接收，会堵住
				//select {
				//case <-ctx.Done():
				//	er = fmt.Errorf("响应没有人接收 %w", ctx.Err())
				//case val <- Resp{
				//	Reply: newReply,
				//	Err:   er,
				//}:
				//}

				// 第一个响应
				select {
				case val <- Resp{
					Reply: newReply,
					Err:   er,
				}:
					// 如果能够成功发送数据到通道 val，则执行这里的代码
				default:
					// 如果通道 val 是满的（对于无缓冲通道）或者发送操作被其他因素阻塞，
				}
				wg.Done()
			}()
		}
		wg.Wait()
		return nil
		//var eg errgroup.Group
		//for _, instance := range instances {
		//	eg.Go(func() error {
		//		iCC, er := grpc.Dial(instance.Addr, c.opts...)
		//		if er != nil {
		//			return er
		//		}
		//		return invoker(ctx, method, req, reply, iCC, opts...)
		//	})
		//}
		//return eg.Wait()
	}
}

type broadcastKey struct {
}

func useBroadcast(ctx context.Context) (context.Context, <-chan Resp) {
	ch := make(chan Resp)
	return context.WithValue(ctx, broadcastKey{}, ch), ch
}

func isBroadcast(ctx context.Context) (chan Resp, bool) {
	val, ok := ctx.Value(broadcastKey{}).(chan Resp)
	return val, ok
}

type Resp struct {
	Err   error
	Reply any
}
