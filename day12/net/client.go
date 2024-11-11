package net

import (
	"encoding/binary"
	"net"
	"time"
)

type Client struct {
	//network any
	//addr string
	//timeout time.Duration
	conn net.Conn
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

func (c *Client) Send(reqData string) (string, error) {
	// 写数据
	lenRep := len(reqData)
	req := make([]byte, lenRep+numOfLengthBytes)
	binary.BigEndian.PutUint64(req[:numOfLengthBytes], uint64(lenRep))
	copy(req[numOfLengthBytes:], reqData)
	_, err := c.conn.Write(req)
	if err != nil {
		return "", err
	}

	// 读数据
	repLenBs := make([]byte, numOfLengthBytes)
	_, err = c.conn.Read(repLenBs)
	if err != nil {
		return "", err
	}
	repLen := binary.BigEndian.Uint64(repLenBs)
	repData := make([]byte, repLen)
	_, err = c.conn.Read(repData)
	if err != nil {
		return "", err
	}

	return string(repData), nil
}
