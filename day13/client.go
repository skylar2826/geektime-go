package rpc

import (
	"context"
	"errors"
	"geektime-go/day13/compressor"
	"geektime-go/day13/compressor/gzip"
	"geektime-go/day13/message"
	"geektime-go/day13/serialize"
	"geektime-go/day13/serialize/json"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"strconv"
	"time"
)

type Client struct {
	pool       pool.Pool
	serialize  serialize.Serializer
	compressor compressor.Compressor
}

type ClientOpts func(c *Client)

func ClientWithSerialize(s serialize.Serializer) ClientOpts {
	return func(c *Client) {
		c.serialize = s
	}
}

func ClientWithCompressor(c compressor.Compressor) ClientOpts {
	return func(client *Client) {
		client.compressor = c
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
		pool:       p,
		serialize:  &json.Serializer{},
		compressor: &gzip.Compressor{},
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
				meta := make(map[string]string, 2)
				if deadline, ok := ctx.Deadline(); ok {
					meta["deadline"] = strconv.FormatInt(deadline.UnixMilli(), 10)
				}
				if isOneWay(ctx) {
					meta["one-way"] = "true"
				}
				// 先序列化再压缩
				if c.compressor.Code() != 0 {
					reqData, err = c.compressor.Compress(reqData)
					if err != nil {
						return []reflect.Value{
							resVal,
							reflect.ValueOf(err),
						}
					}
				}

				req := &message.Request{
					RequestID:   1,
					Version:     2,
					Compressor:  c.compressor.Code(),
					Serializer:  c.serialize.Code(),
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Meta:        meta,
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
				resData := res.Data
				if len(resData) > 0 {
					// 先解压再反序列化
					if res.Compressor != 0 {
						if res.Compressor != c.compressor.Code() {
							return []reflect.Value{
								resVal,
								reflect.ValueOf(errors.New("服务端指定压缩方法客户端不支持")),
							}
						}
						resData, err = c.compressor.UnCompress(resData)
						if err != nil {
							return []reflect.Value{
								resVal,
								reflect.ValueOf(err),
							}
						}
					}
					err = c.serialize.Decode(resData, resVal.Interface())
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
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var resp *message.Response
	var err error
	ch := make(chan struct{})
	go func() {
		resp, err = c.doInvoke(ctx, request)
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-ch:
		return resp, err
	}
}

func (c *Client) doInvoke(ctx context.Context, request *message.Request) (*message.Response, error) {
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
	if isOneWay(ctx) {
		return nil, errors.New("这是单向调用，没有返回值")
	}

	var resBs []byte
	resBs, err = AcceptMsg(conn)
	if err != nil {
		_ = c.pool.Close(conn)
		return nil, err
	}

	return resBs, nil
}
