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
	c.String(http.StatusOK, "index path: %v", c.Path)
}

func main() {
	r := gee.New()

	r.GET("/", indexHandler)

	r.Use(gee.Logger()) // 根Group中间件添加

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.String(http.StatusOK, "url path: %v", c.Path)
		})
		v1.GET("/:name", func(c *gee.Context) {
			fmt.Println(c.Param("name"))
			c.String(http.StatusOK, "fullPath: %s name: %v", c.Path, c.Param("name"))
		})
	}

	v2 := r.Group("/v2")
	{
		v2.Use(gee.Forbidden()) // Abort测试
		v2.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello gee")
		})
	}

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})
	log.Fatal(r.Run(":8888"))
}
