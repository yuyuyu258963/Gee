package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc Used by Gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHttp
type Engine struct {
	router map[string]HandlerFunc
}

// implement ListenAdnServe interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := e.router[joinStr(req.Method, "-", req.URL.Path)]; ok {
		handler(&Context{Request: req, ResponseWriter: w})
	} else {
		fmt.Fprintf(w, "404 Not Found %s", req.URL)
	}
}

// 小写开头package外不可见
// register a handler with method-pattern to gee
func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := joinStr(method, "-", pattern)
	e.router[key] = handler
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

// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{make(map[string]HandlerFunc)}
}
