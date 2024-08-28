package __template_and_file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Context struct {
	W http.ResponseWriter

	RespData       []byte
	RespStatusCode int

	R              *http.Request
	PathParams     map[string]string
	MatchRoute     string // 命中的完整路由
	templateEngine templateEngine

	queryValues url.Values
}

func (c *Context) Render(tagName string, data any) error {
	var err error
	c.RespData, err = c.templateEngine.Render(c.R.Context(), tagName, data)
	if err != nil {
		c.SystemError()
		return err
	}
	c.RespStatusCode = http.StatusOK
	return nil
}

func (c *Context) BadRequest(v ...interface{}) {
	respData := "Bad Request..."
	if len(v) > 0 && v[0] != nil {
		respData = fmt.Sprint(respData, "err: ", v)
	}
	c.WriteJson(http.StatusBadRequest, respData)
}

func (c *Context) SystemError(v ...interface{}) {
	respData := "Internal Server Error..."
	if len(v) > 0 && v[0] != nil {
		respData = fmt.Sprint(respData, "err: ", v)
	}
	c.WriteJson(http.StatusInternalServerError, respData)
}

func (c *Context) RequestOk(v ...interface{}) {
	var respData interface{} = "ok"
	if len(v) > 0 && v[0] != nil {
		respData = v
	}
	c.WriteJson(http.StatusOK, respData)
}

func (c *Context) NotFound() {
	c.WriteJson(http.StatusNotFound, "Not Found")
}

func (c *Context) ReadJson(v interface{}) error {
	data, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) WriteJson(code int, v interface{}) error {
	//c.W.WriteHeader(code)
	c.RespStatusCode = code
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.RespData = data
	//_, err = c.W.Write(data)
	//if err != nil {
	//	return err
	//}
	return nil
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) Reset(w http.ResponseWriter, R *http.Request, templateEngine templateEngine) {
	c.W = w
	c.R = R
	c.PathParams = make(map[string]string, 1)
	c.templateEngine = templateEngine
}

func (c *Context) QueryValue(key string) (string, error) {
	if c.queryValues == nil {
		c.queryValues = c.R.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return vals[0], nil
}

func (c *Context) PathValue(key string) (string, error) {
	path, ok := c.PathParams[key]
	if !ok {
		return "", errors.New(fmt.Sprintf("pathParams key: %s not found", key))
	}
	return path, nil
}
