package rpc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

type Server struct {
	services map[string]Service
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]Service, 16),
	}
}

func (s *Server) registerService(service Service) {
	s.services[service.Name()] = service
}

func (s *Server) Start(network, addr string) error {
	listen, err := net.Listen(network, addr)
	if err != nil {
		// 比较常见的是端口占用
		return err
	}

	for {
		//if errors.Is(err, net.ErrClosed) || err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
		//	return err // 其他err可接受，continue
		//}
		var conn net.Conn
		conn, err = listen.Accept()

		if err != nil {
			return err
		}
		go func() {
			if er := s.handleConn(conn); er != nil {
				_ = conn.Close()
			}
		}()
	}

}

func (s *Server) handleConn(conn net.Conn) error {
	// 读数据
	lenBs := make([]byte, numOfLengthBytes)
	_, err := conn.Read(lenBs)
	if err != nil {
		return err
	}
	// 大顶端 || 小顶端 用客户端协商
	length := binary.BigEndian.Uint64(lenBs)
	reqBs := make([]byte, length)
	_, err = conn.Read(reqBs)
	if err != nil {
		return err
	}

	// 写数据
	var respData []byte
	respData, err = s.handleMsg(reqBs)
	if err != nil {
		return err
	}
	lenRes := len(respData)
	res := make([]byte, lenRes+numOfLengthBytes)
	binary.BigEndian.PutUint64(res[:numOfLengthBytes], uint64(lenRes))
	copy(res[numOfLengthBytes:], respData)
	_, err = conn.Write(res)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleMsg(reqData []byte) ([]byte, error) {
	req := &Request{}
	err := json.Unmarshal(reqData, req)
	if err != nil {
		return nil, err
	}
	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("服务不存在")
	}

	val := reflect.ValueOf(service)
	method := val.MethodByName(req.MethodName)
	in := make([]reflect.Value, 2)
	// 暂时无法从req中获取
	in[0] = reflect.ValueOf(context.Background())

	inReq := reflect.New(method.Type().In(1).Elem())
	err = json.Unmarshal(req.Arg, inReq.Interface())
	if err != nil {
		return nil, err
	}
	in[1] = inReq
	results := method.Call(in)
	// result[0]是返回值
	// result[1]是err
	if results[1].Interface() != nil {
		return nil, results[1].Interface().(error)
	}
	var resp []byte
	resp, err = json.Marshal(results[0].Interface())
	if err != nil {
		return nil, err
	}
	return resp, nil
}
