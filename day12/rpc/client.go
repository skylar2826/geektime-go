package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"time"
)

func InitClientProxy(network, addr string, timeout time.Duration, service Service) error {
	c, err := NewClient(network, addr, timeout)
	if err != nil {
		return err
	}
	return setFuncField(service, c)
}

// 支持远程调用方法
func setFuncField(service Service, proxy proxy) error {
	if service == nil {
		return errors.New("服务不能为空")
	}

	typ := reflect.TypeOf(service)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return errors.New("服务只支持一级指针")
	}

	val := reflect.ValueOf(service).Elem()
	typ = typ.Elem()

	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)

		if fieldVal.CanSet() {
			fnVal := reflect.MakeFunc(fieldTyp.Type, func(args []reflect.Value) (results []reflect.Value) {
				ctx := args[0].Interface().(context.Context)
				resVal := reflect.New(fieldTyp.Type.Out(0).Elem())

				reqData, err := json.Marshal(args[1].Interface())
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}
				req := &Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Arg:         reqData,
				}

				var res *Response
				res, err = proxy.invoke(ctx, req)
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}

				err = json.Unmarshal(res.Data, resVal.Interface())
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}
				return []reflect.Value{
					resVal,
					reflect.Zero(reflect.TypeOf(new(error)).Elem()),
				}
			})
			fieldVal.Set(fnVal)
		}
	}
	return nil
}

type Client struct {
	pool pool.Pool
	ConnMsg
}

func NewClient(network, addr string, timeout time.Duration) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap:  1,
		MaxCap:      30,
		MaxIdle:     10,
		IdleTimeout: time.Minute,
		Factory: func() (interface{}, error) {
			return net.DialTimeout(network, addr, timeout)
		},
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
	})

	if err != nil {
		return nil, err
	}
	return &Client{
		pool: p,
	}, nil

}

func (c *Client) invoke(ctx context.Context, request *Request) (*Response, error) {
	req, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}
	var res []byte
	res, err = c.Send(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := &Response{}
	err = json.Unmarshal(res, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Send(ctx context.Context, data []byte) ([]byte, error) {
	val, err := c.pool.Get()
	conn := val.(net.Conn)
	if err != nil {
		_ = c.pool.Close(conn)
		return nil, err
	}
	err = c.SendMsg(data, conn)
	if err != nil {
		_ = c.pool.Close(conn)
		return nil, err
	}

	var res []byte
	res, err = c.AcceptMsg(conn)
	if err != nil {
		_ = c.pool.Close(conn)
		return nil, err
	}

	return res, nil
}
