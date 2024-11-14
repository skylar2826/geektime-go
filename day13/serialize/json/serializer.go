package json

import "encoding/json"

type Serializer struct {
}

func (s *Serializer) Code() uint8 {
	return 1
}

func (s *Serializer) Encode(val any) ([]byte, error) {
	return json.Marshal(val)
}

func (s *Serializer) Decode(bs []byte, val any) error {
	return json.Unmarshal(bs, val)
}
