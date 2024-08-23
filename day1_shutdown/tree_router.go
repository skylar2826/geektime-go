package __shutdown

import (
	"errors"
	"net/http"
	"sort"
	"strings"
)

var ErrorInvalidRoute = errors.New("invalid router pattern")
var ErrorInvalidMethod = errors.New("invalid method")

type HandlerBasedOnTree struct {
	//root *node
	forest map[string]*node
}

func NewHandlerBasedOnTree() Handler {
	forest := make(map[string]*node, len(supportMethods))
	for _, m := range supportMethods {
		forest[m] = newRootNode(m)
	}
	return &HandlerBasedOnTree{forest}
}

var supportMethods = [4]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
}

func (h *HandlerBasedOnTree) createSubTree(root *node, paths []string, handler handlerFunc) {
	cur := root
	for _, path := range paths {
		nn := createNode(path)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handler
}

// 如果是 * ，则必须是最后一个；
// 非法结构： /abc*, /*/acv
func (h *HandlerBasedOnTree) validatePattern(pattern string) error {
	pos := strings.Index(pattern, "*")

	if pos > 0 {
		// 保证在最后一个
		// 不是 /*/abc结构
		if pos != len(pattern)-1 {
			return ErrorInvalidRoute
		}
		// 保证不是/abc*结构，而是/*结构
		if pattern[pos-1] != '/' {
			return ErrorInvalidRoute
		}
	}
	return nil
}

func (h *HandlerBasedOnTree) Route(method string, pattern string, handler func(c *Context)) error {
	err := h.validatePattern(pattern)
	if err != nil {
		return err
	}

	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	cur, ok := h.forest[method]
	if !ok {
		return ErrorInvalidMethod
	}

	for index, path := range paths {
		matchChild, found := h.findMatchChild(cur, path, nil)
		if found && matchChild.nodeType != nodeTypeAny {
			cur = matchChild
		} else {
			h.createSubTree(cur, paths[index:], handler)
		}
	}

	// 离开循环说明我们加入的短路径
	// 比如：先加入 /order/detail， 再加入/order
	cur.handler = handler
	return nil
}

func (h *HandlerBasedOnTree) findMatchChild(root *node, path string, c *Context) (*node, bool) {
	candidates := make([]*node, 0, 2)
	for _, child := range root.children {
		if child.matchFunc(path, c) {
			candidates = append(candidates, child)
		}
	}

	if len(candidates) == 0 {
		return nil, false
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].nodeType < candidates[j].nodeType
	})

	return candidates[len(candidates)-1], true
}

func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {

	handler, ok := h.findRouter(c.R.Method, c.R.URL.Path, c)
	if ok {
		handler(c)
	} else {
		c.NotFound()
		return
	}
}

func (h *HandlerBasedOnTree) findRouter(method string, pattern string, c *Context) (handlerFunc, bool) {
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	cur, ok := h.forest[method]
	if !ok {
		return nil, false
	}

	for _, path := range paths {
		matchChild, ok := h.findMatchChild(cur, path, c)
		if !ok {
			return nil, false
		}
		cur = matchChild
	}

	if cur.handler == nil {
		return nil, false
	}
	return cur.handler, true
}
