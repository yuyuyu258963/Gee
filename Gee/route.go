package gee

import (
	"fmt"
)

type router struct {
	root *routeTreeNode
}

// 路由树的节点，将path按照"/"进行分割，然后依次添加到路由树上
// 问题是如何区分GET和POST这些请求对应的操作，还是用map?
type routeTreeNode struct {
	nodes    map[string]*routeTreeNode
	handlers map[string][]HandlerFunc
	isLeaf   bool
}

// crete router
func newRouter() *router {
	return &router{root: newRouteTreeNode()}
}

// create RouteTreeNode
func newRouteTreeNode() *routeTreeNode {
	return &routeTreeNode{
		nodes:    make(map[string]*routeTreeNode),
		handlers: make(map[string][]HandlerFunc),
		isLeaf:   false,
	}
}

// addRoute can add a fn to monitor a request with method + pattern
func (r *router) addRoute(method string, pattern string, fn HandlerFunc) {
	urlPath := splitStr(pattern, "?")[0] // 避免query params
	nodeLeaf := r.walkWithCreate(urlPath)
	nodeLeaf.isLeaf = true
	nodeLeaf.addHandler(method, fn)
}

// call handler with Context
func (r *router) handle(c *Context) {
	// 当路由找到对应路径的且是被注册的节点
	if leafNode := r.walk(c.Path); leafNode != nil && leafNode.isLeaf {
		leafNode.callHandler(c.Method, c)
	} else {
		c.Writer.Write([]byte(fmt.Sprintf("404 Not Found %s", c.Path)))
	}
}

// fund the routeTreeNode which is url point to
func (r *router) walk(url string) *routeTreeNode {
	urlItems := splitStr(url, urlSep)
	var root *routeTreeNode = r.root
	var ok bool
	for i := range urlItems {
		if root, ok = root.nodes[urlItems[i]]; !ok {
			return nil
		}
	}
	return root
}

// create route tree with url
// if found a node was not create then create a new node
func (r *router) walkWithCreate(url string) *routeTreeNode {
	urlItems := splitStr(url, "/")
	var root *routeTreeNode = r.root
	var tempRoot *routeTreeNode
	var ok bool
	// fmt.Println(len(urlItems), root)
	for i := range urlItems {
		if tempRoot, ok = root.nodes[urlItems[i]]; ok {
			root = tempRoot
		} else {
			newNode := newRouteTreeNode()
			root.nodes[urlItems[i]] = newNode
			root = newNode
		}
	}
	root.isLeaf = true
	return root
}

// add a handler at routeTreeNode with method and fn
func (r *routeTreeNode) addHandler(method string, fn HandlerFunc) {
	if handlers, ok := r.handlers[method]; ok {
		handlers = append(handlers, fn)
		r.handlers[method] = handlers
	} else {
		r.handlers[method] = []HandlerFunc{fn}
	}
}

// call all handler set on the routeTreeNode
func (r *routeTreeNode) callHandler(method string, c *Context) {
	for i := range r.handlers[method] {
		r.handlers[method][i](c)
	}
}
