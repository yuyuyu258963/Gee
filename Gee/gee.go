package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc Used by Gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHttp
type Engine struct {
	*RouteGroup
	route *router
}

// 因为routeGroup中的操作和Engine上的方法有很多重叠所以写到这个文件中
type RouteGroup struct {
	prefix      string //包含从到当前路由的前缀
	parent      *RouteGroup
	nextGroups  []*RouteGroup
	middlewares []HandlerFunc
	engine      *Engine //因为后序需要直接通过group对象进行操作，但是还需要原始的框架的能力
}

func newEngine() *Engine {
	e := &Engine{route: newRouter()}
	e.RouteGroup = newRouteGroup("", nil, e) //总控的前缀为空字符串
	return e
}

// 创建一个路由分组
func newRouteGroup(prefix string, parent *RouteGroup, e *Engine) *RouteGroup {
	return &RouteGroup{
		prefix:      prefix,
		nextGroups:  make([]*RouteGroup, 0),
		middlewares: make([]HandlerFunc, 0),
		parent:      parent,
		engine:      e,
	}
}

// implement ListenAdnServe interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(req, w) // 每个请求处理的开始时创建一个上下文
	e.route.handle(c)
}

// func (e *Engine) Group(prefix string) (rg *RouteGroup) {
// 	return e.RouteGroup.Group(prefix)
// }

// GER defines a method to add GET request
// func (e *Engine) GET(path string, handler HandlerFunc) {
// 	e.route.addRoute("GET", path, handler)
// }

// // POST defines a method to add POST request
// func (e *Engine) POST(path string, handler HandlerFunc) {
// 	e.route.addRoute("POST", path, handler)
// }

// run and listen request at a port
func (e *Engine) Run(port string) error {
	return http.ListenAndServe(port, e)
}

func (rg *RouteGroup) GET(path string, handler HandlerFunc) {
	rg.engine.route.addRoute("GET", rg.prefix+path, handler)
}

func (rg *RouteGroup) POST(path string, handler HandlerFunc) {
	rg.engine.route.addRoute("POST", rg.prefix+path, handler)
}

// 嵌套地添加分组
func (rg *RouteGroup) Group(prefix string) (g *RouteGroup) {
	prefix = rg.prefix + "/" + prefix

	if prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	g = newRouteGroup(prefix, rg, rg.engine)
	rg.nextGroups = append(rg.nextGroups, g)
	return g
}

func (e *Engine) Test() {
	fmt.Println(e.route)
}

// New is the constructor of gee.Engine
func New() *Engine {
	return newEngine()
}
