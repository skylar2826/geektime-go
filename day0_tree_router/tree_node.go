package __tree_router

import "strings"

const (
	nodeTypeRoot = iota
	nodeTypeAny
	nodeTypeParam
	nodeTypeReg
	nodeTypeStatic
)
const cusAny = "*"

//matchFunc：判断是否匹配，匹配后将必要数据写入Context; 必要数据指路径参数

type matchFunc func(path string, c *Context) bool

type node struct {
	pattern   string
	nodeType  int
	children  []*node
	handler   handlerFunc
	matchFunc matchFunc
}

func newStaticNode(path string) *node {
	return &node{
		pattern:  path,
		nodeType: nodeTypeStatic,
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			return path == p && p != "*"
		},
	}
}

// *
func newAnyNode() *node {
	return &node{
		pattern:  cusAny,
		nodeType: nodeTypeAny,
		// 我们不允许通配符后有节点(/*/friend)，所以没有children
		matchFunc: func(p string, c *Context) bool { return true },
	}
}

func newRootNode(method string) *node {
	return &node{
		pattern:   method,
		nodeType:  nodeTypeRoot,
		children:  make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool { panic("never call me") },
	}
}

// /user/:id
func newParamNode(path string) *node {
	paramName := path[1:]
	return &node{
		pattern:  path,
		nodeType: nodeTypeParam,
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			if c != nil {
				c.PathParams[paramName] = p
			}
			// 不支持： /user/:id/*
			return p != cusAny
		},
	}
}

func createNode(path string) *node {
	if path == "*" {
		return newAnyNode()
	}
	if strings.HasPrefix(path, ":") {
		return newParamNode(path)
	}
	return newStaticNode(path)
}
