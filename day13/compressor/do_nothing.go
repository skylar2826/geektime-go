package compressor

type doNothingCompressor struct {
}

func (d doNothingCompressor) Code() uint8 {
	return 0
}

func (d doNothingCompressor) Compress(data []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (d doNothingCompressor) UnCompress(data []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
