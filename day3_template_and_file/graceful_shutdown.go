package __template_and_file

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var ErrorRequestTimeout = errors.New("handle request timeout")

type GracefulShutdown struct {
	reqCnt     int64
	closing    int32
	zeroReqCnt chan struct{}
}

func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		zeroReqCnt: make(chan struct{}),
	}
}

func (g *GracefulShutdown) ShutdownFilter(next Filter) Filter {
	return func(c *Context) {
		closing := atomic.LoadInt32(&g.closing)
		if closing > 0 {
			c.RespStatusCode = http.StatusServiceUnavailable
			//c.W.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		atomic.AddInt64(&g.reqCnt, 1)
		next(c)
		reqCnt := atomic.AddInt64(&g.reqCnt, -1)

		if closing > 0 && reqCnt == 0 {
			g.zeroReqCnt <- struct{}{}
		}
	}
}

func (g *GracefulShutdown) RejectNewRequestAndWaiting(ctx context.Context) error {
	atomic.AddInt32(&g.closing, 1)

	if atomic.LoadInt64(&g.reqCnt) == 0 {
		return nil
	}

	done := ctx.Done()
	select {
	case <-done:
		fmt.Printf("接口处理超时")
		return ErrorRequestTimeout
	case <-g.zeroReqCnt:
		fmt.Printf("接口处理完成")
		return nil
	}
}

func WaitForShutdown(hooks ...Hook) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutdownSignals...)

	select {
	case sig := <-signals:
		fmt.Printf("get signal %s, application will shutdown\n", sig)
		time.AfterFunc(time.Second*10, func() {
			fmt.Println("请求未处理完成, 直接关闭")
			os.Exit(-1)
		})

		for _, h := range hooks {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			err := h(ctx)
			if err != nil {
				fmt.Printf("hook err: %v\n", err)
				cancel()
			}
		}
		os.Exit(0)
	}
}
