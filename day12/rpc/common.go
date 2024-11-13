package rpc

import (
	"encoding/binary"
	"net"
)

var numOfDataLength = 8

type ConnMsg struct {
}

func (c *ConnMsg) SendMsg(data []byte, conn net.Conn) error {
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

func (c *ConnMsg) AcceptMsg(conn net.Conn) ([]byte, error) {
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
