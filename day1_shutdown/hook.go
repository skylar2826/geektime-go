package __shutdown

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrorHookTimeout = errors.New("hook timeout")

type Hook func(ctx context.Context) error

func BuildCloseServerHook(servers ...Server) Hook {
	return func(ctx context.Context) error {
		wg := sync.WaitGroup{}
		doneCh := make(chan struct{})
		wg.Add(len(servers))

		for _, server := range servers {
			go func(svr Server) {
				defer wg.Done()
				err := svr.Shutdown(ctx)
				if err != nil {
					fmt.Printf("shutdown server error: %v\n", err)
				}
			}(server)
		}

		go func() {
			wg.Wait()
			doneCh <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			fmt.Printf("closing servers timeout\n")
			return ErrorHookTimeout
		case <-doneCh:
			fmt.Printf("closing servers done\n")
			return nil
		}
	}
}
