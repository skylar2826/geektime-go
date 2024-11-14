package message

import (
	"encoding/binary"
)

type Response struct {
	HeadLength uint32 `json:"HeadLength"`
	BodyLength uint32 `json:"BodyLength"`
	RequestID  uint32 `json:"RequestID"`
	Version    uint8  `json:"Version"`
	Compressor uint8  `json:"Compressor"`
	Serializer uint8  `json:"Serializer"`

	Error []byte `json:"Error"`

	Data []byte `json:"data"`
}

func (r *Response) CalculateHeaderLength() {
	r.HeadLength = uint32(15 + len(r.Error))
}

func (r *Response) CalculateBodyLength() {
	r.BodyLength = uint32(len(r.Data))
}

func EncodeResp(resp *Response) []byte {
	respBs := make([]byte, resp.HeadLength+resp.BodyLength)
	binary.BigEndian.PutUint32(respBs[:4], resp.HeadLength)
	binary.BigEndian.PutUint32(respBs[4:8], resp.BodyLength)
	binary.BigEndian.PutUint32(respBs[8:12], resp.RequestID)
	respBs[12] = resp.Version
	respBs[13] = resp.Compressor
	respBs[14] = resp.Serializer

	copy(respBs[15:len(resp.Error)+15], resp.Error)
	copy(respBs[resp.HeadLength:], resp.Data)

	return respBs
}

func DecodeResp(respBs []byte) *Response {
	resp := &Response{}

	resp.HeadLength = binary.BigEndian.Uint32(respBs[:4])
	resp.BodyLength = binary.BigEndian.Uint32(respBs[4:8])
	resp.RequestID = binary.BigEndian.Uint32(respBs[8:12])
	resp.Version = respBs[12]
	resp.Compressor = respBs[13]
	resp.Serializer = respBs[14]

	if resp.HeadLength > 15 {
		resp.Error = respBs[15:resp.HeadLength]
	}

	if resp.BodyLength != 0 {
		resp.Data = respBs[resp.HeadLength:]
	}

	return resp
}
