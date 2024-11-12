package rpc

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_setFuncField(t *testing.T) {
	testCases := []struct {
		name    string
		service Service
		mock    func(ctrl *gomock.Controller) Proxy
		wantErr error
	}{
		{
			name:    "nil",
			service: nil,
			mock: func(ctrl *gomock.Controller) Proxy {
				return NewMockProxy(ctrl)
			},
			wantErr: errors.New("service 不允许为 nil"),
		},
		{
			name:    "pointer",
			service: UserService{},
			mock: func(ctrl *gomock.Controller) Proxy {
				return NewMockProxy(ctrl)
			},
			wantErr: errors.New("只支持指向结构体的一级指针"),
		},
		{
			name:    "user service",
			service: &UserService{},
			mock: func(ctrl *gomock.Controller) Proxy {
				p := NewMockProxy(ctrl)

				p.EXPECT().invoke(gomock.Any(), &Request{
					ServiceName: "user-service",
					MethodName:  "GetById",
					Arg:         []byte(`{"Id": "123"}`),
				}).Return(
					&Response{}, nil)
				return p
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			err := setFuncField(tc.service, tc.mock(ctrl))
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			var resp *GetByIdResp
			resp, err = tc.service.(*UserService).GetById(context.Background(), &GetByIdReq{Id: "123"})
			assert.Equal(t, err, tc.wantErr)
			t.Log(resp)
		})
	}
}

// mockgen -destination=day12/rpc/mock_proxy_gen_test.go -package=rpc -source=day12/rpc/types.go Proxy
