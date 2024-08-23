package __filter_builder

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

func (h *HandlerBasedOnTree) createSubTree(root *node, paths []string, handler func(c *Context), builders ...FilterBuilder) {
	cur := root
	for _, path := range paths {
		nn := createNode(path)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	// server.Route 注册handler
	if handler != nil {
		cur.handler = handler
	}
	// server.Use 注册filter
	if builders != nil {
		cur.mdls = builders
	}

	cur.matchRoute = strings.Join(paths, "/")
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

// 注册路由
func (h *HandlerBasedOnTree) Route(method string, pattern string, handler func(c *Context), builders ...FilterBuilder) error {
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

	if pattern != "" { // 不是根节点
		for index, path := range paths {
			matchChild, found, _ := h.findMatchChild(cur, path, nil)
			if found && matchChild.nodeType != nodeTypeAny {
				cur = matchChild
			} else {
				h.createSubTree(cur, paths[index:], handler, builders...)
				return nil
			}
		}
	}

	// 离开循环说明我们加入的短路径
	// 比如：先加入 /order/detail， 再加入/order
	if handler != nil {
		cur.handler = handler
	}
	if builders != nil {
		cur.mdls = builders
	}
	cur.matchRoute = pattern
	return nil
}

func (h *HandlerBasedOnTree) findMatchChild(root *node, path string, c *Context) (*node, bool, []*node) {
	candidates := make([]*node, 0, 2)
	for _, child := range root.children {
		if child.matchFunc(path, c) {
			candidates = append(candidates, child)
		}
	}

	if len(candidates) == 0 {
		return nil, false, make([]*node, 0)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].nodeType < candidates[j].nodeType
	})

	return candidates[len(candidates)-1], true, candidates
}

type matchInfo struct {
	node *node
	mdls []FilterBuilder
}

func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	mi, ok := h.findRouter(c.R.Method, c.R.URL.Path, c)

	if ok {

		c.MatchRoute = mi.node.matchRoute

		var rootHandler Filter = func(c *Context) {
			if mi.node.handler != nil {
				mi.node.handler(c)
			}
		}
		// 从后往前包，靠近的先执行
		for i := len(mi.mdls) - 1; i >= 0; i-- {
			b := mi.mdls[i]
			rootHandler = b(rootHandler)
		}

		rootHandler(c)
	} else {
		c.NotFound()
		return
	}
}

func (h *HandlerBasedOnTree) findRouter(method string, pattern string, c *Context) (*matchInfo, bool) {
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	root, ok := h.forest[method]
	if !ok {
		return nil, false
	}

	cur := root
	for _, path := range paths {
		matchChild, ok, _ := h.findMatchChild(cur, path, c)
		if !ok {
			return nil, false
		}

		cur = matchChild
	}

	if cur.handler == nil {
		return nil, false
	}

	mi := &matchInfo{
		node: cur,
		mdls: h.findFilterBuilders(root, paths...),
	}

	return mi, true
}

func (h *HandlerBasedOnTree) findFilterBuilders(root *node, paths ...string) []FilterBuilder {
	queue := []*node{root}
	res := make([]FilterBuilder, 0, 16)

	// /a/b 从根节点开始这条路径上的所有filterBuilder都要执行，从后向前执行
	for _, path := range paths {
		var children []*node
		for _, child := range queue {
			if len(child.mdls) > 0 {
				res = append(res, child.mdls...)
			}
			_, found, matchNodes := h.findMatchChild(child, path, nil)
			if found {
				children = matchNodes
			}
		}
		queue = children
	}

	for _, child := range queue {
		if len(child.mdls) > 0 {
			res = append(res, child.mdls...)
		}
	}
	return res
}

func (h *HandlerBasedOnTree) AddFilterBuilders(method string, pattern string, builders ...FilterBuilder) error {
	return h.Route(method, pattern, nil, builders...)
}
