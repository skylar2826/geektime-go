package rpc

import (
	"context"
	"errors"
	"geektime-go/day13/message"
	"geektime-go/day13/serialize"
	"geektime-go/day13/serialize/json"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"time"
)

type Client struct {
	pool      pool.Pool
	serialize serialize.Serializer
}

type ClientOpts func(c *Client)

func ClientWithSerialize(s serialize.Serializer) ClientOpts {
	return func(c *Client) {
		c.serialize = s
	}
}

func NewClient(network, addr string, timeout time.Duration, opts ...ClientOpts) (*Client, error) {
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
	c := &Client{
		pool:      p,
		serialize: &json.Serializer{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Client) InitService(service Service) error {
	return c.setFuncField(service, c)
}

// 支持远程调用方法
func (c *Client) setFuncField(service Service, proxy proxy) error {
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

				reqData, err := c.serialize.Encode(args[1].Interface())
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}
				req := &message.Request{
					RequestID:   1,
					Version:     2,
					Compressor:  3,
					Serializer:  c.serialize.Code(),
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Meta:        map[string]string{},
					Data:        reqData,
				}
				req.CalculateHeaderLength()
				req.CalculateBodyLength()

				var res *message.Response
				res, err = proxy.invoke(ctx, req)
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}

				var resErr error
				if len(res.Error) > 0 {
					resErr = errors.New(string(res.Error))
				}
				if len(res.Data) > 0 {
					err = c.serialize.Decode(res.Data, resVal.Interface())
					if err != nil {
						return []reflect.Value{
							resVal,
							reflect.ValueOf(err),
						}
					}
				}

				var resErrVal reflect.Value
				if resErr == nil {
					resErrVal = reflect.Zero(reflect.TypeOf(new(error)).Elem())
				} else {
					resErrVal = reflect.ValueOf(resErr)
				}

				return []reflect.Value{
					resVal,
					resErrVal,
				}
			})
			fieldVal.Set(fnVal)
		}
	}
	return nil
}

func (c *Client) invoke(ctx context.Context, request *message.Request) (*message.Response, error) {
	req := message.EncodeReq(request)
	resBs, err := c.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := message.DecodeResp(resBs)
	return resp, nil
}

func (c *Client) Send(ctx context.Context, data []byte) ([]byte, error) {
	val, err := c.pool.Get()
	conn := val.(net.Conn)
	if err != nil {
		_ = c.pool.Close(conn)
		return nil, err
	}

	_, err = conn.Write(data)
	if err != nil {
		_ = c.pool.Close(conn)
		return nil, err
	}

	var resBs []byte
	resBs, err = AcceptMsg(conn)
	if err != nil {
		_ = c.pool.Close(conn)
		return nil, err
	}

	return resBs, nil
}
