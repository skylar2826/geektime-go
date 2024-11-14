package rpc

import (
	"encoding/binary"
	"net"
)

var numOfDataLength = 8

func AcceptMsg(conn net.Conn) ([]byte, error) {
	resLenBs := make([]byte, numOfDataLength)
	_, err := conn.Read(resLenBs)
	if err != nil {
		return nil, err
	}
	headerLength := binary.BigEndian.Uint32(resLenBs[:4])
	bodyLength := binary.BigEndian.Uint32(resLenBs[4:])
	resLen := headerLength + bodyLength
	res := make([]byte, resLen)

	copy(res[:8], resLenBs)
	_, err = conn.Read(res[8:])
	if err != nil {
		return nil, err
	}
	return res, nil
}
