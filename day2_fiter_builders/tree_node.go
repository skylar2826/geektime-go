package __filter_builder

import (
	"regexp"
	"strings"
)

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
	pattern  string
	nodeType int
	children []*node
	handler  HandlerFunc
	//handler    Filter
	matchFunc  matchFunc
	matchRoute string
	mdls       []FilterBuilder
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

// 正则节点 /user/:id(/[1-9]/)
func newRegNode(paramsName string, pattern string, path string) *node {
	return &node{
		pattern:  path,
		nodeType: nodeTypeReg,
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return false
			}

			if c != nil {
				c.PathParams[paramsName] = p
			}
			return re.MatchString(p)
		},
	}
}

func createNode(path string) *node {
	if path == "*" {
		return newAnyNode()
	}
	if strings.HasPrefix(path, ":") {
		re, err := regexp.Compile(`:(\w+)\(([^)]*)\)`)
		if err != nil {
			panic("创建noe 正则失败")
		}
		matches := re.FindStringSubmatch(path)

		if matches == nil {
			return newParamNode(path)
		} else {
			return newRegNode(matches[1], matches[2], path)
		}
	}
	return newStaticNode(path)
}
