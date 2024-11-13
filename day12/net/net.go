package net

import (
	"context"
	"encoding/binary"
	"net"
	"time"
)

var numOfDataLength = 8

type connMsg struct {
}

func (c *connMsg) SendMsg(ctx context.Context, data []byte, conn net.Conn) error {
	lenData := len(data)
	req := make([]byte, lenData+numOfDataLength)
	binary.BigEndian.PutUint64(req[:numOfDataLength], uint64(lenData))
	copy(req[numOfDataLength:], data)
	_, err := conn.Write(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *connMsg) acceptMsg(ctx context.Context, conn net.Conn) ([]byte, error) {
	resLenBs := make([]byte, numOfDataLength)
	_, err := conn.Read(resLenBs)
	if err != nil {
		return nil, err
	}
	resLen := binary.BigEndian.Uint64(resLenBs)

	res := make([]byte, resLen)
	_, err = conn.Read(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Client struct {
	conn net.Conn
	connMsg
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

func (c *Client) Send(ctx context.Context, data []byte) ([]byte, error) {
	err := c.SendMsg(ctx, data, c.conn)
	if err != nil {
		_ = c.conn.Close()
		return nil, err
	}

	var res []byte
	res, err = c.acceptMsg(ctx, c.conn)
	if err != nil {
		_ = c.conn.Close()
		return nil, err
	}

	return res, nil
}

type Server struct {
	network string
	addr    string
	connMsg
}

func NewServer(network string, addr string) (*Server, error) {
	return &Server{
		network: network,
		addr:    addr,
	}, nil
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
		req, err = s.acceptMsg(context.Background(), conn)
		if err != nil {
			_ = conn.Close()
			return err
		}
		var res []byte
		res, err = s.handleService(req)
		if err != nil {
			// ? 业务出错，需要包装返回？
			_ = conn.Close()

			return err
		}
		err = s.SendMsg(context.Background(), res, conn)
		if err != nil {
			_ = conn.Close()

			return err
		}
	}
}

func (s *Server) handleService(req []byte) ([]byte, error) {
	res := append(req, []byte(" response")...)
	return res, nil
}
