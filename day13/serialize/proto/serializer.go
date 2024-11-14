package proto

import (
	"errors"
	"github.com/golang/protobuf/proto"
)

type Serializer struct {
}

func (s *Serializer) Code() uint8 {
	return 2
}

func (s *Serializer) Encode(val any) ([]byte, error) {
	msg, ok := val.(proto.Message)
	if !ok {
		return nil, errors.New("proto:val不是message类型")
	}
	return proto.Marshal(msg)
}

func (s *Serializer) Decode(bs []byte, val any) error {
	msg, ok := val.(proto.Message)
	if !ok {
		return errors.New("proto:val不是message类型")
	}
	return proto.Unmarshal(bs, msg)
}
