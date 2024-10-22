package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context with Req information and ResponseWriter
type Context struct {
	// http origin
	Req    *http.Request
	Writer http.ResponseWriter
	// request info
	Method string
	Path   string
	Params map[string]string

	// Controller
	index   int           // 当前在处理的控制函数
	handles []HandlerFunc // ...中间件，注册的请求处理函数

	// response info
	StatusCode int
}

type H map[string]interface{}

// create a new Context with request and ResponseWriter
func newContext(req *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		Req:     req,
		Writer:  w,
		Method:  req.Method,
		Path:    req.URL.Path,
		index:   -1,
		handles: make([]HandlerFunc, 0),
	}
}

// 往后执行处理函数
func (c *Context) Next() {
	for c.index = c.index + 1; c.index < len(c.handles); c.index++ {
		c.handles[c.index](c)
	}
}

// 终止执行后续的处理函数
func (c *Context) Abort() {
	c.index = len(c.handles)
}

// modify code to status
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// add key : val to response header
func (c *Context) setHeader(key string, val string) {
	c.Writer.Header().Set(key, val)
}

// 获得请求path
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Query enable get data by key
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// QueryWithDefault while not found the key
func (c *Context) QueryWithDefault(key string, defaultValue string) string {
	val := c.Req.URL.Query().Get(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

// PostForm get the value
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// PostFormWithDefault find key mapped val
// if not found the key then return defaultValue
func (c *Context) PostFormWithDefault(key string, defaultValue string) string {
	val := c.Req.FormValue(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

// Data can write bytes type of data to response
func (c *Context) Data(status int, data []byte) {
	c.Status(status)
	c.Writer.Write(data)
}

// String can Write string to Response with format
func (c *Context) String(status int, format string, values ...interface{}) {
	c.setHeader("Content-Type", "text/plain")
	c.Status(status)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON can write a map[string]interface{} typed data
func (c *Context) JSON(status int, h H) {
	c.setHeader("Content-Type", "application/json")
	c.Status(status)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(h); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// HTML can write a html code to body
func (c *Context) HTML(status int, htm string) {
	c.setHeader("Content-Type", "text/html")
	c.Status(status)
	c.Writer.Write([]byte(htm))
}
