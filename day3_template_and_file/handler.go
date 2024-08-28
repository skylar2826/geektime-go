package __template_and_file

type HandlerFunc func(c *Context)

type Handler interface {
	//http.Handler // 不使用原因：该方法将context拆开
	Route(method string, pattern string, handler func(c *Context), builders ...FilterBuilder) error
	ServeHTTP(c *Context)
	AddFilterBuilders(method string, pattern string, builders ...FilterBuilder) error
}
