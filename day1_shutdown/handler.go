package __shutdown

type handlerFunc func(c *Context)

type Handler interface {
	//http.Handler // 不使用原因：该方法将context拆开
	Routable
	ServeHTTP(c *Context)
}
