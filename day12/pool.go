package day12

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

type Pool struct {
	// 空闲队列
	idlesConns chan *idleConn
	// 请求队列
	reqConn []connReq

	// 最大连接数
	maxCnt int
	// 当前连接数
	cnt int
	// 初始连接数
	initCnt int

	// 最大连接时间
	maxIdleTime time.Duration

	factory func() (net.Conn, error)

	lock sync.Mutex
}

type idleConn struct {
	c              net.Conn
	lastActiveTime time.Time
}

type connReq struct {
	connChan chan net.Conn
}

func NewPool(maxCnt int, maxIdleCnt int, maxIdleTime time.Duration, initCnt int, factory func() (net.Conn, error)) (*Pool, error) {
	idleConns := make(chan *idleConn, maxIdleCnt)

	if initCnt > maxIdleCnt {
		return nil, errors.New("初始连接数不能大于最大连接数")
	}

	for i := 0; i < initCnt; i++ {
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		idleConns <- &idleConn{c: conn, lastActiveTime: time.Now()}
	}

	res := &Pool{
		idlesConns:  idleConns,
		maxCnt:      maxCnt,
		cnt:         0,
		maxIdleTime: maxIdleTime,
		factory:     factory,
	}

	return res, nil
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}

	for {
		select {
		case ic := <-p.idlesConns:
			// 拿到空闲连接
			// 判断是否超时
			if ic.lastActiveTime.Add(p.maxIdleTime).Before(time.Now()) {
				_ = ic.c.Close()
				continue
			}
			return ic.c, nil
		default:
			p.lock.Lock()
			if p.cnt >= p.maxCnt {
				req := connReq{connChan: make(chan net.Conn, 1)}
				p.reqConn = append(p.reqConn, req)
				p.lock.Unlock()
				// 没有空闲连接
				select {
				case <-ctx.Done():
					// 超时
					// 选项1：从req中删除自己（不行，没得删）
					// 选项2：转发
					go func() {
						c := <-req.connChan
						_ = p.Put(ctx, c)
					}()
					return nil, ctx.Err()
				case c := <-req.connChan:
					// 等别人归还
					return c, nil
				}
			} else {
				// 有空闲连接
				c, err := p.factory()
				if err != nil {
					return nil, err
				}
				p.cnt++
				p.lock.Unlock()
				return c, nil
			}
		}

	}

}

func (p *Pool) Put(ctx context.Context, conn net.Conn) error {
	p.lock.Lock()
	if len(p.reqConn) > 0 {
		// 取队尾元素，队首元素等待时间长，容易过期
		req := p.reqConn[len(p.reqConn)-1]
		p.reqConn = p.reqConn[:(len(p.reqConn) - 1)]
		p.lock.Unlock()
		req.connChan <- conn
		return nil
	}
	p.lock.Unlock()
	ic := &idleConn{
		c:              conn,
		lastActiveTime: time.Now(),
	}
	select {
	case p.idlesConns <- ic:
	default:
		// 空闲队列满了
		_ = conn.Close()
		p.cnt--
	}
	return nil
}
