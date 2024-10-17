package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc Used by Gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHttp
type Engine struct {
	route *router
}

func newEngine() *Engine {
	return &Engine{route: newRouter()}
}

// implement ListenAdnServe interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(req, w) // 每个请求处理的开始时创建一个上下文
	e.route.handle(c)
}

// GER defines a method to add GET request
func (e *Engine) GET(path string, handler HandlerFunc) {
	e.route.addRoute("GET", path, handler)
}

// POST defines a method to add POST request
func (e *Engine) POST(path string, handler HandlerFunc) {
	e.route.addRoute("POST", path, handler)
}

// run and listen request at a port
func (e *Engine) Run(port string) error {
	return http.ListenAndServe(port, e)
}

func (e *Engine) Test() {
	fmt.Println(e.route)
}

// New is the constructor of gee.Engine
func New() *Engine {
	return newEngine()
}
