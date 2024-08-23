package __tree_router

type HandlerBasedOnMap struct {
	handlers map[string]func(c *Context)
}

func (h *HandlerBasedOnMap) key(method string, pattern string) string {
	return method + "#" + pattern
}

func (h *HandlerBasedOnMap) Route(method string, pattern string, handler func(c *Context)) error {
	key := h.key(method, pattern)
	h.handlers[key] = handler
	return nil
}

func (h *HandlerBasedOnMap) ServeHTTP(c *Context) {
	key := h.key(c.R.Method, c.R.URL.Path)
	ctx := NewContext(c.W, c.R)

	if handler, ok := h.handlers[key]; ok {
		handler(ctx)
	} else {
		ctx.NotFound()
	}
}

// 确保HandlerBasedOnMap实现Handler; 一般加在实现类之上
var _ Handler = &HandlerBasedOnMap{}

// 不像暴露创建细节，所以提供New方法
func NewHandlerBasedOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: make(map[string]func(c *Context)),
	}
}
