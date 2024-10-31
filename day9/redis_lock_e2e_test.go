package day9

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisE2ELock_Lock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	testCases := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		wantErr    error
		expiration time.Duration
		key        string
		wantLock   *Lock
	}{
		{
			name: "key exist test", // 锁被他人抢走了
			key:  "key1",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				res, err := rdb.Set(ctx, "key1", "value1", time.Minute).Result()
				require.NoError(t, err)
				assert.Equal(t, res, "OK")
			},
			expiration: time.Minute,
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				//ERR unknown command `getdel`, with args beginning with: `key1`,
				//res, err := rdb.GetDel(ctx, "key1").Result()
				res, err := rdb.Get(ctx, "key1").Result()
				require.NoError(t, err)
				assert.Equal(t, res, "value1")
				rdb.Del(ctx, "key1")
			},
			wantErr: ErrFailedToPreemptLock,
		},
		{
			name:       "locked",
			key:        "key3",
			before:     func(t *testing.T) {},
			expiration: time.Minute,
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				res, err := rdb.Get(ctx, "key3").Result()
				require.NoError(t, err)
				assert.NotEmpty(t, res)

				rdb.Del(ctx, "key3")
			},
			wantLock: &Lock{
				key: "key3",
			},
		},
		{
			name:       "locked with preempt timeout",
			key:        "key2",
			expiration: time.Minute,
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				res, err := rdb.Set(ctx, "key2", "value2", time.Second).Result()
				require.NoError(t, err)
				assert.Equal(t, res, "OK")
				time.Sleep(time.Second * 2) // 等待key1过期
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				_, err := rdb.Get(ctx, "key2").Result()
				require.NoError(t, err)
				rdb.Del(ctx, "key2")
			},
			wantLock: &Lock{
				key: "key2",
			},
		},
	}

	client := NewRedisLockClient(rdb)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			tc.before(t)
			l, err := client.tryLock(ctx, tc.key, tc.expiration)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				tc.after(t)
				return
			}

			assert.Equal(t, l.key, tc.wantLock.key)
			assert.NotEmpty(t, l.value)
			assert.NotNil(t, l.client)
			tc.after(t)
		})
	}
}

func TestRedisE2ELock_Unlock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	testCases := []struct {
		name    string
		key     string
		before  func(t *testing.T)
		after   func(t *testing.T)
		wantErr error
		lock    *Lock
	}{
		{
			name: "lock not exist",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			lock: &Lock{
				key:    "unlock_key1",
				value:  "123",
				client: rdb,
			},
			wantErr: ErrLockNotExist,
		},
		{
			name: "lock preempt",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				_, err := rdb.Set(ctx, "unlock_key2", "value1", time.Second*10).Result()
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				res, err := rdb.Get(ctx, "unlock_key2").Result()
				require.NoError(t, err)
				assert.Equal(t, res, "value1")
				rdb.Del(ctx, "unlock_key2")
			},
			lock: &Lock{
				key:        "unlock_key2",
				value:      "123",
				client:     rdb,
				expiration: time.Minute,
			},
			wantErr: ErrLockNotExist,
		},
		{
			name: "unlocked",
			key:  "unlock_key3",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				_, err := rdb.Set(ctx, "unlock_key3", "123", time.Minute).Result()
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				res, err := rdb.Exists(ctx, "unlock_key3").Result()
				require.NoError(t, err)
				assert.Equal(t, res, int64(0))
			},
			lock: &Lock{
				key:    "unlock_key3",
				value:  "123",
				client: rdb,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			tc.before(t)
			err := tc.lock.Unlock(ctx)
			assert.Equal(t, err, tc.wantErr)
			tc.after(t)
		})
	}
}

func TestRedisE2ELock_Refresh(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	testCases := []struct {
		name    string
		key     string
		before  func(t *testing.T)
		after   func(t *testing.T)
		wantErr error
		lock    *Lock
	}{
		{
			name: "lock not exist",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			lock: &Lock{
				key:    "Refresh_key1",
				value:  "123",
				client: rdb,
			},
			wantErr: ErrFailedLock,
		},
		{
			name: "lock preempt",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				_, err := rdb.Set(ctx, "Refresh_key2", "value1", time.Second*10).Result()
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				timeout, err := rdb.TTL(ctx, "Refresh_key2").Result()
				require.NoError(t, err)
				// 刷新成功，timeout 约等于 1分钟；没刷新成功，timeout < 10s
				assert.True(t, timeout <= time.Second*10)
				res, err := rdb.Get(ctx, "Refresh_key2").Result()
				require.NoError(t, err)
				assert.Equal(t, res, "value1")
				rdb.Del(ctx, "Refresh_key2")
			},
			lock: &Lock{
				key:        "Refresh_key2",
				value:      "123",
				client:     rdb,
				expiration: time.Minute,
			},
			wantErr: ErrFailedLock,
		},
		{
			name: "Refreshed",
			key:  "Refresh_key3",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				_, err := rdb.Set(ctx, "Refresh_key3", "123", time.Second*10).Result()
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				timeout, err := rdb.TTL(ctx, "Refresh_key3").Result()
				require.NoError(t, err)
				assert.True(t, timeout > time.Second*50)
				var res int64
				res, err = rdb.Exists(ctx, "Refresh_key3").Result()
				require.NoError(t, err)
				assert.Equal(t, res, int64(1))
				rdb.Del(ctx, "Refresh_key3")
			},
			lock: &Lock{
				key:        "Refresh_key3",
				value:      "123",
				client:     rdb,
				expiration: time.Minute,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			tc.before(t)
			err := tc.lock.Refresh(ctx)
			assert.Equal(t, err, tc.wantErr)
			tc.after(t)
		})
	}
}

func ExampleLock_Refresh() {
	var l *Lock
	// 间隔多久续约一次?
	// 刷新超时时间怎么设置？
	// 出现错误了怎么办？
	//		停止服务
	//		retry
	stopChan := make(chan struct{})
	errChan := make(chan error)
	timeoutChan := make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		cnt := 0
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)
				cnt = 0
				cancel()
				if errors.Is(err, context.DeadlineExceeded) {
					timeoutChan <- struct{}{}
					continue
				}
				if err != nil {
					errChan <- err
					close(stopChan)
					close(errChan)
					close(timeoutChan)
					return
				}

			case <-stopChan:
				err := l.Unlock(context.Background())
				errChan <- err
				close(stopChan)
				close(errChan)
				close(timeoutChan)
				return
			case <-timeoutChan:
				if cnt > 10 {
					return
				}
				cnt++
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)
				cancel()
				if errors.Is(err, context.DeadlineExceeded) {
					timeoutChan <- struct{}{}
					continue
				}
				if err != nil {
					errChan <- err
					close(stopChan)
					close(errChan)
					close(timeoutChan)
					return
				}
			}

		}

	}()

	// 没有循环的业务：
	select {
	case <-errChan:
		// 续约失败，去中断业务
		return
	default:

	}

	// 业务结束，进行stop
	stopChan <- struct{}{}
}

func ExampleLock_AutoRefresh() {
	var l *Lock
	errChan := make(chan error, 1)
	go func() {
		err := l.AutoRefresh(time.Second, time.Second*10)
		errChan <- err
	}()

	select {
	case <-errChan:
		// 续约失败，去中断业务
		return
	default:

	}
}
