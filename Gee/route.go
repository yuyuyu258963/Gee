package gee

import (
	"net/http"
	"strings"
)

// 路由树实现的方式
// - root: 通过map将不同的请求方法映射到不同的请求分支上
// - handles：通过 "method-pattern" 记录路由树上所有的请求处理方法
type router struct {
	root    map[string]*node
	handles map[string]HandlerFunc
}

// 路由树的节点，将path按照"/"进行分割，然后依次添加到路由树上
// 问题是如何区分GET和POST这些请求对应的操作，还是用map?
// TODO 是否要考虑将wildChildren 和普通的节点进行区分
type node struct {
	part     string // 当前节点对应的路由路径名
	pattern  string // 叶节点标记注册的路由路径
	children []*node
	isWild   bool // 是否为精确匹配，part含有*时为true
}

// crete new router with init
func newRouter() *router {
	return &router{root: make(map[string]*node),
		handles: make(map[string]HandlerFunc)}
}

// 创建一个路由树节点
func newNode() *node {
	return &node{
		children: make([]*node, 0),
		isWild:   false,
	}
}

// 往路由树上添加一个节点
func (r *router) addRoute(method string, pattern string, fn HandlerFunc) {
	parts := parsePattern(pattern)
	if _, ok := r.root[method]; !ok {
		r.root[method] = newNode()
	}
	r.root[method].insert(pattern, parts, 0)
	key := joinStr(method, "-", pattern)
	r.handles[key] = fn
}

// 获取根据method和pattern请求得到最后的路由树的叶节点
func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)

	root, ok := r.root[method]
	if !ok {
		return nil, nil
	}
	// 存储动态路由的匹配
	params := make(map[string]string)
	leafNode := root.search(pattern, searchParts, 0)

	if leafNode != nil { // 只有匹配上的时候才去解析上面的Param
		parts := parsePattern(leafNode.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return leafNode, params
	}
	return nil, nil
}

// 根据当前的请求在路由树上查找处理函数
func (r *router) getHandle(c *Context) HandlerFunc {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		return r.handles[key]
	}
	return notFoundHandle
}

// 匹配成功的节点，用于插入
// 优先匹配出非通配的节点
func (n *node) mathChild(part string) *node {
	var chNode *node
	for _, child := range n.children {
		if child.part == part {
			return child
		}
		if child.isWild {
			chNode = child
		}
	}
	return chNode
}

// 匹配出所有的子节点，
// 先匹配的是非通配节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part && !child.isWild {
			nodes = append(nodes, child)
		}
	}
	for _, child := range n.children {
		if child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 递归的向树中添加节点，若未发现这么一条路径则创建
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果当前的节点已经到了最后一个，或者当前的节点是*通配节点
	if height == len(parts) {
		n.pattern = pattern // 标记为叶节点，只有叶节点含有不为空的pattern
		return
	}

	part := parts[height]
	child := n.mathChild(part)
	if child == nil {
		child = &node{children: make([]*node, 0),
			part:   part,
			isWild: part[0] == ':' || part[0] == '*', // 两种通配匹配的方案
		}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height+1)
}

// 递归的在路由树上查找节点
// - 优先匹配非通配节点
func (n *node) search(pattern string, parts []string, height int) *node {
	//找到了最后或者是当前part含有通配符*的时候
	// 如果匹配结束了，或者匹配到了通配符*就检查是不是被注册过，没注册够就返回nil
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children { // 如果没有在上面找到
		if next := child.search(pattern, parts, height+1); next != nil {
			return next
		}
	}
	return nil
}

// 未找到页面的HandlerFunc
func notFoundHandle(c *Context) {
	c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
}
