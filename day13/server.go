package rpc

import (
	"context"
	"errors"
	"geektime-go/day13/compressor"
	"geektime-go/day13/compressor/gzip"
	"geektime-go/day13/message"
	"geektime-go/day13/serialize"
	"geektime-go/day13/serialize/json"
	"net"
	"reflect"
	"strconv"
	"time"
)

type Server struct {
	network    string
	addr       string
	services   map[string]reflectionStub
	serializes map[uint8]serialize.Serializer
	compressor map[uint8]compressor.Compressor
}

func NewServer(network string, addr string) (*Server, error) {
	s := &Server{
		network:    network,
		addr:       addr,
		services:   make(map[string]reflectionStub, 16),
		serializes: make(map[uint8]serialize.Serializer, 4),
		compressor: make(map[uint8]compressor.Compressor, 4),
	}

	j := &json.Serializer{}
	c := &gzip.Compressor{}
	s.serializes[j.Code()] = j
	s.compressor[c.Code()] = c

	return s, nil
}

func (s *Server) registerService(service Service) {
	s.services[service.Name()] = reflectionStub{
		service:    service,
		value:      reflect.ValueOf(service),
		serialize:  s.serializes,
		compressor: s.compressor,
	}
}

func (s *Server) registerSerialize(serializer serialize.Serializer) {
	s.serializes[serializer.Code()] = serializer
}

func (s *Server) registerCompressor(compressor compressor.Compressor) {
	s.compressor[compressor.Code()] = compressor
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

		req := message.DecodeReq(reqBs)

		var res *message.Response
		res, err = s.invoke(req)
		if err != nil {
			if res == nil {
				res = &message.Response{}
			}
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

func (s *Server) invoke(req *message.Request) (*message.Response, error) {
	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("服务不存在")
	}
	ctx := context.Background()
	var deadlineStr string
	cancel := func() {}
	if deadlineStr, ok = req.Meta["deadline"]; ok {
		if deadline, er := strconv.ParseInt(deadlineStr, 10, 64); er == nil {
			ctx, cancel = context.WithDeadline(ctx, time.UnixMilli(deadline))
		}
	}
	oneWay, exist := req.Meta["one-way"]
	if exist && oneWay == "true" {
		go func() {
			_, _ = service.invoke(ctx, req)
			cancel()
		}()
		return nil, errors.New("这是单向调用，没有返回值")
	} else {
		resData, err := service.invoke(ctx, req)
		cancel()

		resp := &message.Response{
			RequestID:  req.RequestID,
			Version:    req.Version,
			Compressor: req.Compressor,
			Serializer: req.Serializer,
			Data:       resData,
		}

		if err != nil {
			if resp == nil {
				resp = &message.Response{}
			}
			resp.Error = []byte(err.Error())
		}

		return resp, nil
	}
}

type reflectionStub struct {
	service    Service
	value      reflect.Value
	serialize  map[uint8]serialize.Serializer
	compressor map[uint8]compressor.Compressor
}

func (r *reflectionStub) invoke(ctx context.Context, req *message.Request) ([]byte, error) {
	// 先解压再反序列化
	var reqData []byte
	var c compressor.Compressor
	if req.Compressor != 0 {
		var ok bool
		c, ok = r.compressor[req.Compressor]
		if !ok {
			return nil, errors.New("客户端指定的压缩算法服务端不存在")
		}
		var err error
		reqData, err = c.UnCompress(req.Data)
		if err != nil {
			return nil, err
		}
	}

	serializer, ok := r.serialize[req.Serializer]
	if !ok {
		return nil, errors.New("序列化协议不存在")
	}

	serviceElem := reflect.ValueOf(r.service)
	method := serviceElem.MethodByName(req.MethodName)

	in := make([]reflect.Value, 2)

	in[0] = reflect.ValueOf(ctx)

	inReq := reflect.New(method.Type().In(1).Elem())
	err := serializer.Decode(reqData, inReq.Interface())
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
		res, er = serializer.Encode(result[0].Interface())
		if er != nil {
			return nil, er
		}
		// 先序列化再压缩
		if req.Compressor != 0 {
			c, ok = r.compressor[req.Compressor]
			if !ok {
				return nil, errors.New("客户端指定的压缩算法服务端不存在")
			}
			res, er = c.Compress(res)
			if er != nil {
				return nil, er
			}
		}
	}

	return res, result[1].Interface().(error)
}
