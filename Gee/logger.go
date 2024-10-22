package gee

import (
	"log"
	"net/http"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		c.Next() // call other handlers
		// call back
		log.Printf("[%d] %s int %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func Forbidden() HandlerFunc {
	return func(c *Context) {
		c.String(http.StatusForbidden, "Forbidden")
		c.Abort()
	}
}
