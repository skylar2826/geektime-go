package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRespEncodeDecode(t *testing.T) {
	testCases := []struct {
		name string
		resp *Response
	}{
		{
			name: "normal",
			resp: &Response{
				RequestID:  123,
				Version:    1,
				Compressor: 13,
				Serializer: 14,
				Error:      []byte("error"),
				Data:       []byte("hello world"),
			},
		},

		{
			name: "no error",
			resp: &Response{
				RequestID:  123,
				Version:    1,
				Compressor: 13,
				Serializer: 14,
				Data:       []byte("hello world"),
			},
		},
		{
			name: "data has \n",
			resp: &Response{
				//HeadLength: 10,
				//BodyLength: 11,
				RequestID:  123,
				Version:    1,
				Compressor: 13,
				Serializer: 14,
				Error:      []byte("error"),
				Data:       []byte("hello \n world"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.resp.CalculateHeaderLength()
			tc.resp.CalculateBodyLength()
			respBs := EncodeResp(tc.resp)
			req := DecodeResp(respBs)
			assert.Equal(t, tc.resp, req)
		})
	}
}
