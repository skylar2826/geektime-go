package message

import (
	"bytes"
	"encoding/binary"
)

var separator byte = '\n'
var metaSeparator byte = '\r'

type Request struct {
	HeadLength uint32 `json:"HeadLength"`
	BodyLength uint32 `json:"BodyLength"`
	RequestID  uint32 `json:"RequestID"`
	Version    uint8  `json:"Version"`
	Compressor uint8  `json:"Compressor"`
	Serializer uint8  `json:"Serializer"`

	ServiceName string `json:"serviceName"`
	MethodName  string `json:"methodName"`

	Meta map[string]string `json:"Meta"`

	Data []byte `json:"arg"`
}

func (r *Request) CalculateHeaderLength() {
	lenMeta := 0
	for key, value := range r.Meta {
		lenMeta += len(key) + len(value) + 2
	}

	r.HeadLength = uint32(15 + len(r.ServiceName) + 1 + len(r.MethodName) + 1 + lenMeta)
}

func (r *Request) CalculateBodyLength() {
	r.BodyLength = uint32(len(r.Data))
}

func EncodeReq(req *Request) []byte {
	reqBs := make([]byte, req.HeadLength+req.BodyLength)
	binary.BigEndian.PutUint32(reqBs[:4], req.HeadLength)
	binary.BigEndian.PutUint32(reqBs[4:8], req.BodyLength)
	binary.BigEndian.PutUint32(reqBs[8:12], req.RequestID)
	reqBs[12] = req.Version
	reqBs[13] = req.Compressor
	reqBs[14] = req.Serializer

	cur := reqBs[15:]
	copy(cur[:len(req.ServiceName)], req.ServiceName)
	cur = cur[len(req.ServiceName):]
	cur[0] = separator
	copy(cur[1:], req.MethodName)
	cur = cur[(len(req.MethodName) + 1):]
	cur[0] = separator
	for key, value := range req.Meta {
		copy(cur[1:len(key)+1], key)
		cur[len(key)+1] = metaSeparator
		copy(cur[len(key)+2:], value)
		cur[len(key)+2+len(value)] = separator
		cur = cur[len(key)+2+len(value):]
	}

	copy(reqBs[req.HeadLength:], req.Data)

	return reqBs
}

func DecodeReq(reqBs []byte) *Request {
	req := &Request{}

	req.HeadLength = binary.BigEndian.Uint32(reqBs[:4])
	req.BodyLength = binary.BigEndian.Uint32(reqBs[4:8])
	req.RequestID = binary.BigEndian.Uint32(reqBs[8:12])
	req.Version = reqBs[12]
	req.Compressor = reqBs[13]
	req.Serializer = reqBs[14]

	header := reqBs[15:req.HeadLength]
	index := bytes.IndexByte(header, separator)
	if index != -1 {
		req.ServiceName = string(header[:index])
	}
	header = header[index+1:]
	index = bytes.IndexByte(header, separator)
	if index != -1 {
		req.MethodName = string(header[:index])
	}

	header = header[index+1:]
	index = bytes.IndexByte(header, separator)
	meta := make(map[string]string, 16)
	for {
		if index == -1 {
			break
		}
		i := bytes.IndexByte(header[:index], metaSeparator)
		key := header[:i]
		value := header[i+1 : index]
		meta[string(key)] = string(value)
		header = header[index+1:]
		index = bytes.IndexByte(header, separator)
	}

	req.Meta = meta

	if req.BodyLength != 0 {
		req.Data = reqBs[req.HeadLength:]
	}

	return req
}
