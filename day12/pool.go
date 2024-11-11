package day12

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

type Pool struct {
	idleConns   chan *idleConn // 空闲队列 chan的容量 =》 最大空闲数量
	reqQueue    []connReq      // 等待队列，不使用chan原因：chan有最大容量，超过容量的会被进一步阻塞
	maxCnt      int            // 最大连接数
	cnt         int            // 当前连接数
	maxIdleTime time.Duration  // 最大超时时间
	factory     func() (net.Conn, error)
	lock        sync.Mutex
}

func NewPool(initCnt int, maxIdleConns int, maxCnt int, maxIdleTime time.Duration, factory func() (net.Conn, error)) (*Pool, error) {
	if initCnt > maxCnt {
		return nil, errors.New("初始连接数大于最大连接数")
	}

	idleConns := make(chan *idleConn, maxIdleConns)
	for i := 0; i <= initCnt; i++ {
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		idleConns <- &idleConn{
			conn:           conn,
			lastActiveTime: time.Now(),
		}
	}

	res := &Pool{
		idleConns:   idleConns,
		maxCnt:      maxCnt,
		maxIdleTime: maxIdleTime,
		factory:     factory,
	}
	return res, nil
}

type idleConn struct {
	conn           net.Conn
	lastActiveTime time.Time
}

type connReq struct {
	conn chan net.Conn
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}

	for {
		select {
		case i := <-p.idleConns:
			// 存在空闲连接
			// 检查是否超时
			if i.lastActiveTime.Add(p.maxIdleTime).Before(time.Now()) {
				// 超时，获取下一个空闲连接
				continue
			}
			return i.conn, nil
		default:
			// 不存在空闲连接
			p.lock.Lock()

			if p.cnt >= p.maxCnt {
				// 超过最大数量
				cReq := connReq{
					conn: make(chan net.Conn, 1),
				}

				p.reqQueue = append(p.reqQueue, cReq)

				select {
				case c := <-cReq.conn:
					// 等到别人归还
					p.lock.Unlock()
					return c, nil
				case <-ctx.Done():
					// 转发该空闲连接
					go func() {
						c := <-cReq.conn
						_ = p.Put(context.Background(), c)
					}()
					p.lock.Unlock()
					return nil, ctx.Err()
				}
			}

			// 没有超过最大连接数 =》 创建连接
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

func (p *Pool) Put(ctx context.Context, conn net.Conn) error {
	if len(p.reqQueue) > 0 {
		// 有阻塞请求，唤醒阻塞请求
		p.lock.Lock()
		reqLen := len(p.reqQueue)
		// 取队尾元素，队尾元素比队首元素后放入，大概率没有超期
		req := p.reqQueue[reqLen-1]
		p.reqQueue = p.reqQueue[:reqLen-1]
		p.lock.Unlock()
		req.conn <- conn
		return nil
	}

	p.lock.Unlock()

	idleC := &idleConn{
		conn: conn,
	}
	select {
	case p.idleConns <- idleC:
	default:
		// 空闲队列满了
		_ = conn.Close()
		p.lock.Lock()
		p.cnt--
		p.lock.Unlock()
	}

	return nil
}
