package day9

import (
	"context"
	"geektime-go/day9/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisLock_Lock(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) redis.Cmdable
		key        string
		expiration time.Duration
		wantErr    error
	}{
		{
			name:       "setNX err",
			key:        "key1",
			expiration: time.Second * 2,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewBoolResult(false, context.DeadlineExceeded)
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Second*2).Return(status)
				return cmd
			},
			wantErr: ErrFailedToPreemptLock,
		},
		{
			name:       "not preempt",
			key:        "key1",
			expiration: time.Second * 2,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(false, nil)
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Second*2).Return(res)
				return cmd
			},
			wantErr: ErrFailedToPreemptLock,
		},
		{
			name:       "not preempt",
			key:        "key1",
			expiration: time.Second * 2,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(true, nil)
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Second*2).Return(res)
				return cmd
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			r := NewRedisLockClient(tc.mock(ctrl))
			l, err := r.tryLock(context.Background(), tc.key, tc.expiration)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, l.key, tc.key)
		})
	}
}

func TestRedisLock_Unlock(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		key     string
		value   string
		wantErr error
	}{
		{
			name:  "error",
			key:   "key1",
			value: "value1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().Eval(context.Background(), luaUnlock, []string{"key1"}, "value1").Return(res)
				return cmd
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name:  "error",
			key:   "key1",
			value: "value1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(0)
				cmd.EXPECT().Eval(context.Background(), luaUnlock, []string{"key1"}, "value1").Return(res)
				return cmd
			},
			wantErr: ErrLockNotExist,
		},
		{
			name:  "error",
			key:   "key1",
			value: "value1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(1)
				cmd.EXPECT().Eval(context.Background(), luaUnlock, []string{"key1"}, "value1").Return(res)
				return cmd
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			l := &Lock{
				key:    tc.key,
				value:  tc.value,
				client: tc.mock(ctrl),
			}
			err := l.Unlock(context.Background())
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}

		})
	}
}
