package __tree_router

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	PathParams map[string]string
}

func (c *Context) BadRequest() {
	c.WriteJson(http.StatusBadRequest, "Bad Request")
}

func (c *Context) SystemError() {
	c.WriteJson(http.StatusInternalServerError, "Internal Server Error")
}

func (c *Context) RequestOk() {
	c.WriteJson(http.StatusOK, "OK")
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
	c.W.WriteHeader(code)
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = c.W.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{W: w, R: r, PathParams: make(map[string]string, 1)}
}
