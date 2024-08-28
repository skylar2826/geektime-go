package __template_and_file

import (
	"fmt"
	"net/http"
)

type RequestJson struct {
	Email    string `json:"u_email"`
	Password string `json:"u_password"`
}

type ResponseJson struct {
	Data string `json:"data"`
	Code int    `json:"code"`
}

func SignUp(c *Context) {
	requestJson := &RequestJson{}
	err := c.ReadJson(requestJson)
	if err != nil {
		c.BadRequest()
		return
	}
	fmt.Println("打印参数（mock 验证身份）：", requestJson)
	responseJson := &ResponseJson{
		Data: "hello world",
	}
	err = c.WriteJson(http.StatusOK, responseJson)
	if err != nil {
		c.SystemError()
		return
	}
}
