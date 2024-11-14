package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"geektime-go/day13/message"
	"net"
	"reflect"
)

type Server struct {
	network  string
	addr     string
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
		reqBs, err = AcceptMsg(conn)
		if err != nil {
			_ = conn.Close()
			return err
		}

		// 其实需要从请求中拿到 ctx
		ctx := context.Background()
		req := message.DecodeReq(reqBs)
		var res *message.Response
		res, err = s.invoke(ctx, req)
		if err != nil {
			res.Error = []byte(err.Error())
		}
		res.CalculateHeaderLength()
		res.CalculateBodyLength()
		resBs := message.EncodeResp(res)
		_, err = conn.Write(resBs)
		if err != nil {
			_ = conn.Close()
			return err
		}
	}
}

func (s *Server) invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("服务不存在")
	}
	resData, err := service.invoke(ctx, req.MethodName, req.Data)
	if err != nil {
		return nil, err
	}

	resp := &message.Response{
		RequestID:  req.RequestID,
		Version:    req.Version,
		Compressor: req.Compressor,
		Serializer: req.Serializer,
		Data:       resData,
	}

	return resp, nil
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

	var res []byte
	if result[0].IsNil() {
		return nil, err
	} else {
		var er error
		res, er = json.Marshal(result[0].Interface())
		if er != nil {
			return nil, er
		}
	}

	return res, nil
}
