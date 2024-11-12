package rpc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"time"
)

// InitClientProxy 要为 GetById 之类的函数类型字段赋值
func InitClientProxy(addr string, service Service) error {
	client, err := NewClient(addr)
	if err != nil {
		return err
	}
	return setFuncField(service, client)
}

func setFuncField(service Service, proxy Proxy) error {
	if service == nil {
		return errors.New("service 不允许为 nil")
	}
	val := reflect.ValueOf(service)
	typ := val.Type()
	// 只支持一级指针结构
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return errors.New("只支持指向结构体的一级指针")
	}

	val = val.Elem()
	typ = typ.Elem()

	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)

		if fieldVal.CanSet() {
			fnVal := reflect.MakeFunc(fieldTyp.Type, func(args []reflect.Value) (results []reflect.Value) {

				ctx := args[0].Interface().(context.Context)
				retVal := reflect.New(fieldTyp.Type.Out(0).Elem()) // 返回值的第0个

				reqData, err := json.Marshal(args[1].Interface())
				if err != nil {
					return []reflect.Value{
						retVal,
						reflect.ValueOf(err),
					}
				}
				req := &Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,

					//args[0]是context, args[1] 是参数
					//Args: slice.Map[reflect.Value, any](args, func(idx int, src reflect.Value) any {
					//	return src.Interface()
					//}),
					Arg: reqData,
				}

				var resp *Response
				resp, err = proxy.invoke(ctx, req)

				if err != nil {
					return []reflect.Value{
						retVal,
						reflect.ValueOf(err),
					}
				}

				err = json.Unmarshal(resp.data, retVal.Interface())
				if err != nil {
					return []reflect.Value{
						retVal,
						reflect.ValueOf(err),
					}
				}
				return []reflect.Value{
					retVal,
					reflect.Zero(reflect.TypeOf(new(error)).Elem()),
				}
			})
			fieldVal.Set(fnVal)
		}
	}

	return nil
}

const numOfLengthBytes = 8

type Client struct {
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*3)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) invoke(ctx context.Context, req *Request) (*Response, error) {
	// 发送请求至服务端

	// 转换成二进制数据
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 用连接发送出去
	var resp []byte
	resp, err = c.Send(data)
	if err != nil {
		return nil, err
	}
	return &Response{
		data: resp,
	}, nil
}

func (c *Client) Send(reqData []byte) ([]byte, error) {
	// 写数据
	lenRep := len(reqData)
	req := make([]byte, lenRep+numOfLengthBytes)
	binary.BigEndian.PutUint64(req[:numOfLengthBytes], uint64(lenRep))
	copy(req[numOfLengthBytes:], reqData)
	_, err := c.conn.Write(req)
	if err != nil {
		return nil, err
	}

	// 读数据
	repLenBs := make([]byte, numOfLengthBytes)
	_, err = c.conn.Read(repLenBs)
	if err != nil {
		return nil, err
	}
	repLen := binary.BigEndian.Uint64(repLenBs)
	repData := make([]byte, repLen)
	_, err = c.conn.Read(repData)
	if err != nil {
		return nil, err
	}

	return repData, nil
}
