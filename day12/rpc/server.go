package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

type Server struct {
	network string
	addr    string
	ConnMsg
	services map[string]reflectionStub
}

func NewServer(network string, addr string) (*Server, error) {
	s := &Server{
		network:  network,
		addr:     addr,
		services: make(map[string]reflectionStub, 16),
	}

	return s, nil
}

func (s *Server) registerService(service Service) {
	s.services[service.Name()] = reflectionStub{
		service: service,
		value:   reflect.ValueOf(service),
	}
}

func (s *Server) Start() error {
	listen, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}

	for {
		var conn net.Conn
		conn, err = listen.Accept()
		if err != nil {
			return err
		}

		var reqBs []byte
		reqBs, err = s.AcceptMsg(conn)
		if err != nil {
			_ = conn.Close()
			return err
		}

		// 其实需要从请求中拿到 ctx
		ctx := context.Background()
		req := &Request{}
		err = json.Unmarshal(reqBs, req)
		if err != nil {
			return err
		}

		var res *Response
		res, err = s.invoke(ctx, req)
		if err != nil {
			// ? 业务出错，需要包装返回？
			_ = conn.Close()
			return err
		}

		var resBs []byte
		resBs, err = json.Marshal(res)
		if err != nil {
			return err
		}

		err = s.SendMsg(resBs, conn)
		if err != nil {
			_ = conn.Close()
			return err
		}
	}
}

func (s *Server) invoke(ctx context.Context, req *Request) (*Response, error) {
	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("服务不存在")
	}
	resData, err := service.invoke(ctx, req.MethodName, req.Arg)
	if err != nil {
		return nil, err
	}
	return &Response{
		Data: resData,
	}, nil
}

type reflectionStub struct {
	service Service
	value   reflect.Value
}

func (r *reflectionStub) invoke(ctx context.Context, methodName string, arg []byte) ([]byte, error) {
	serviceElem := reflect.ValueOf(r.service)
	method := serviceElem.MethodByName(methodName)

	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(ctx)

	inReq := reflect.New(method.Type().In(1).Elem())
	err := json.Unmarshal(arg, inReq.Interface())
	if err != nil {
		return nil, err
	}

	in[1] = inReq
	result := method.Call(in)
	if result[1].Interface() != nil {
		return nil, result[1].Interface().(error)
	}
	var data []byte
	data, err = json.Marshal(result[0].Interface())
	if err != nil {
		return nil, err
	}

	return data, nil
}
