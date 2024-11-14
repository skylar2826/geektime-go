package serialize

type Serializer interface {
	Code() uint8
	Encode(val any) ([]byte, error)
	Decode(bs []byte, val any) error
}
