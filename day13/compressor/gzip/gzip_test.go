package gzip

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGZip(t *testing.T) {
	c := &Compressor{}
	str := []byte("hello world!")
	res, err := c.Compress(str)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("加密：", res)
	res, err = c.UnCompress(res)
	if err != nil {
		t.Log(err)
		return
	}
	assert.Equal(t, res, str)
}
