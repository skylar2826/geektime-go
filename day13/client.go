package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"geektime-go/day13/message"
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
				req := &message.Request{
					RequestID:   1,
					Version:     2,
					Compressor:  3,
					Serializer:  4,
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
					err = json.Unmarshal(res.Data, resVal.Interface())
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

type Client struct {
	pool pool.Pool
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
