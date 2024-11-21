package gzip

import (
	"bytes"
	"compress/gzip"
	"errors"
	"golang.org/x/sys/windows/registry"
	"io"
)

type Compressor struct {
}

func (G Compressor) Code() uint8 {
	return 1
}

func (G Compressor) Compress(data []byte) ([]byte, error) {
	res := bytes.NewBuffer(nil)
	w := gzip.NewWriter(res)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

func (G Compressor) UnCompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Close()
	}()

	var res []byte
	res, err = io.ReadAll(r)
	if err != nil && err != io.EOF && !errors.Is(err, registry.ErrUnexpectedType) {
		return nil, err
	}
	return res, nil
}
