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
	services map[string]Service
}

func NewServer(network string, addr string) (*Server, error) {
	s := &Server{
		network:  network,
		addr:     addr,
		services: make(map[string]Service, 16),
	}

	return s, nil
}

func (s *Server) registerService(service Service) {
	s.services[service.Name()] = service
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
		var req []byte

		req, err = s.AcceptMsg(conn)
		if err != nil {
			_ = conn.Close()
			return err
		}

		var res []byte
		// 其实需要从请求中拿到 ctx
		ctx := context.Background()
		res, err = s.handleService(ctx, req)
		if err != nil {
			// ? 业务出错，需要包装返回？
			_ = conn.Close()

			return err
		}
		err = s.SendMsg(res, conn)
		if err != nil {
			_ = conn.Close()

			return err
		}
	}
}

func (s *Server) handleService(ctx context.Context, reqBs []byte) ([]byte, error) {
	var req Request
	err := json.Unmarshal(reqBs, &req)
	if err != nil {
		return nil, err
	}
	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("服务不存在")
	}

	serviceElem := reflect.ValueOf(service)
	method := serviceElem.MethodByName(req.MethodName)

	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(ctx)

	inReq := reflect.New(method.Type().In(1).Elem())
	err = json.Unmarshal(req.Arg, inReq.Interface())
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

	res := &Response{
		Data: data,
	}
	var resBs []byte
	resBs, err = json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return resBs, nil
}
