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
