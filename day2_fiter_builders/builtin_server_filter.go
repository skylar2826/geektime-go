package __filter_builder

import (
	"fmt"
)

// HandleResp 提取resp管理，允许其他中间件修改
func HandleResp(next Filter) Filter {
	return func(c *Context) {
		next(c)
		if c.RespStatusCode != 0 {
			c.W.WriteHeader(c.RespStatusCode)
		}
		_, err := c.W.Write(c.RespData)
		if err != nil {
			fmt.Println("write response error:", err)
		}
	}
}
