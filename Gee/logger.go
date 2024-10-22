package gee

import (
	"log"
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
