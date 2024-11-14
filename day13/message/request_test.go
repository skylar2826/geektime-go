package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
	}{
		{
			name: "normal",
			req: &Request{
				//HeadLength: 10,
				//BodyLength: 11,
				RequestID:   123,
				Version:     1,
				Compressor:  13,
				Serializer:  14,
				ServiceName: "user_service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"trace_id": "555",
					"a/b":      "test",
					"c/d":      "test2",
				},
				Data: []byte("hello world"),
			},
		},

		{
			name: "empty meta",
			req: &Request{
				//HeadLength: 10,
				//BodyLength: 11,
				RequestID:   123,
				Version:     1,
				Compressor:  13,
				Serializer:  14,
				ServiceName: "user_service",
				MethodName:  "GetById",
				Meta:        map[string]string{},
				Data:        []byte("hello world"),
			},
		},
		{
			name: "data has \n",
			req: &Request{
				//HeadLength: 10,
				//BodyLength: 11,
				RequestID:   123,
				Version:     1,
				Compressor:  13,
				Serializer:  14,
				ServiceName: "user_service",
				MethodName:  "GetById",
				Meta:        map[string]string{},
				Data:        []byte("hello \n world"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.req.CalculateHeaderLength()
			tc.req.CalculateBodyLength()
			reqBs := encodeReq(tc.req)
			req := DecodeReq(reqBs)
			assert.Equal(t, tc.req, req)
		})
	}
}
