package day9

import (
	"context"
	"errors"
	"geektime-go/day9/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCache_Set(t *testing.T) {
	testCases := []struct {
		name       string
		key        string
		value      any
		expiration time.Duration
		wantErr    error
		mock       func(ctrl *gomock.Controller) redis.Cmdable
	}{
		{
			name:       "set key",
			key:        "key1",
			value:      "value1",
			expiration: time.Second,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetVal("OK")
				cmd.EXPECT().Set(context.Background(), "key1", "value1", time.Second).Return(status)
				return cmd
			},
		},
		{
			name:       "timeout",
			key:        "key1",
			value:      "value1",
			expiration: time.Second,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().Set(context.Background(), "key1", "value1", time.Second).Return(status)
				return cmd
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name:       "unexpected key",
			key:        "unexpected key",
			value:      "value1",
			expiration: time.Second,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetVal("Not ok")
				cmd.EXPECT().Set(context.Background(), "unexpected key", "value1", time.Second).Return(status)
				return cmd
			},
			wantErr: errors.New(ErrFailedKey),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewRedisCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.key, tc.value, tc.expiration)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

func TestRedisCache_Get(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		wantErr error
		wantVal any
	}{
		{
			name:    "get value",
			key:     "key1",
			wantVal: "value1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStringCmd(context.Background())
				status.SetVal("value1")
				cmd.EXPECT().Get(context.Background(), "key1").Return(status)
				return cmd
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewRedisCache(tc.mock(ctrl))
			val, err := c.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, val, tc.wantVal)
		})
	}
}
