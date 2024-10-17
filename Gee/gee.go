package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc Used by Gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHttp
type Engine struct {
	router *roteTreeRoot
}

// Init
func (e *Engine) Init() {
	e.router = &roteTreeRoot{}
	e.router.Init()
}

// implement ListenAdnServe interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 当路由找到对应路径的且是被注册的节点
	if leafNode := e.router.walk(req.URL.Path); leafNode != nil && leafNode.isLeaf {
		c := &Context{Request: req, ResponseWriter: w}
		leafNode.callHandler(req.Method, c)
	} else {
		fmt.Fprintf(w, "404 Not Found %s", req.URL)
	}
}

// 小写开头package外不可见
// register a handler with method-pattern to gee
func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	urlPath := splitStr(pattern, "?")[0] // 避免query params
	roteTreeLeaf := e.router.walkWithCreate(urlPath)
	roteTreeLeaf.addHandler(method, handler)
}

// GER defines a method to add GET request
func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRoute("GET", path, handler)
}

// POST defines a method to add POST request
func (e *Engine) POST(path string, handler HandlerFunc) {
	e.addRoute("POST", path, handler)
}

// run and listen request at a port
func (e *Engine) Run(port string) error {
	return http.ListenAndServe(port, e)
}

func (e *Engine) Test() {
	fmt.Println(e.router)
}

// New is the constructor of gee.Engine
func New() *Engine {
	e := &Engine{}
	e.Init()
	return e
}
