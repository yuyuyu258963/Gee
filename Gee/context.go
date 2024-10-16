package gee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Context with Request information and ResponseWriter
type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

type H map[string]interface{}

func (c *Context) writeStatus(status int) {
	c.ResponseWriter.WriteHeader(status)
}

// Query enable get data by key
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// QueryWithDefault while not found the key
func (c *Context) QueryWithDefault(key string, defaultValue string) string {
	val := c.Request.URL.Query().Get(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

// PostForm get the value
func (c *Context) PostForm(key string) string {
	return c.Request.PostFormValue(key)
}

// PostFormWithDefault find key mapped val
// if not found the key then return defaultValue
func (c *Context) PostFormWithDefault(key string, defaultValue string) string {
	val := c.Request.PostFormValue(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

// String can Write string to Response with format
func (c *Context) String(status int, s string, format ...any) {
	c.writeStatus(status)
	fmt.Fprintf(c.ResponseWriter, s, format...)
}

// JSON can write a map[string]interface{} typed data
func (c *Context) JSON(status int, h H) {
	c.writeStatus(status)
	jsonStr, err := json.Marshal(h)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %v", err)
	}
	fmt.Fprint(c.ResponseWriter, string(jsonStr))
}

// HTML can write a html code to body
func (c *Context) HTML(status int, htm string) {
	c.writeStatus(status)
	fmt.Fprint(c.ResponseWriter, htm)
}
