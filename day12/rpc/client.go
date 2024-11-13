package rpc

import (
	"context"
	"encoding/json"
	"errors"
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
	conn net.Conn
	ConnMsg
}

func NewClient(network, addr string, timeout time.Duration) (*Client, error) {
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
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
	err := c.SendMsg(data, c.conn)
	if err != nil {
		_ = c.conn.Close()
		return nil, err
	}

	var res []byte
	res, err = c.AcceptMsg(c.conn)
	if err != nil {
		_ = c.conn.Close()
		return nil, err
	}

	return res, nil
}
