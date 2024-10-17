package main

import (
	"fmt"
	gee "gee/Gee"
	"log"
	"net/http"
)

func indexHandler(c *gee.Context) {
	// expect /?name=abc&age=11&cccc=b
	fmt.Println(c.Query("name"))
	fmt.Println(c.Query("age"))
	fmt.Println(c.Query("cccc"))
	c.String(http.StatusOK, "index path: %v", c.Request.URL.Path)
}

func htmHandler(c *gee.Context) {
	c.String(http.StatusOK, `<h1>Hello Gee</h1>
														<p2>html!</p2>`)
}

func htmaHandler(c *gee.Context) {
	c.String(http.StatusOK, `<h1>Hello Gee</h1>
														<p2>html a!</p2>`)
}

func htmabcHandler(c *gee.Context) {
	c.String(http.StatusOK, `<h1>Hello Gee</h1>
														<p2>html abc!</p2>`)
}

func helloHandler(c *gee.Context) {
	c.JSON(http.StatusOK, gee.H{
		"name": c.PostForm("name"),
		"pwd":  c.PostFormWithDefault("pwd", "null"),
	})
}

func main() {
	r := gee.New()
	r.GET("/", indexHandler)
	r.GET("/htm", htmHandler)
	r.GET("/htm/a/c", htmaHandler)
	r.GET("/htm/a/b/c", htmaHandler)
	r.POST("/login", helloHandler)
	log.Fatal(r.Run(":8888"))
}
