package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
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

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	// skip = 3
	// 是因为第0个调用者是Callers本身，第一个是上一层的trace，再上面一层是defer func
	// 这三个函数调用信息对错误回溯没有作用，所以跳过
	n := runtime.Callers(3, pcs[:]) // skip first caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s \n\n", trace(message))
				c.String(http.StatusInternalServerError, "Internal Server Error")
				c.Abort()
			}
		}()

		// use recursion
		c.Next() // all panic cause by other(after) functions and middlers can be recovered
	}
}
