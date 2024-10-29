package day9

import (
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	p := sync.Pool{
		New: func() any {
			t.Log("创建资源")
			return "hello"
		},
	}
	str := p.Get()
	t.Log(str)
	p.Put(str)
	str = p.Get()
	defer p.Put(str)
	t.Log(str)
}
