package net

import (
	"encoding/binary"
	"net"
)

// 长度字段使用的字节数量
const numOfLengthBytes = 8

type Server struct{}

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
	respData := handleMsg(reqBs)
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

func handleMsg(msg []byte) []byte {
	res := make([]byte, len(msg)*2)
	copy(res[:len(msg)], msg)
	copy(res[len(msg):], msg)
	return res
}
