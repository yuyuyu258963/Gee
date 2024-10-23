package main

import (
	"fmt"
	gee "gee/Gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

func indexHandler(c *gee.Context) {
	// expect /?name=abc&age=11&cccc=b
	fmt.Println(c.Query("name"))
	fmt.Println(c.Query("age"))
	fmt.Println(c.Query("cccc"))
	c.String(http.StatusOK, "index path: %v", c.Path)
}

// 模版中可使用的函数
// 例如：
// <p> {{FormatAsDate .Now}} </p>
// 可以在模版中调用该函数
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func serve() {
	r := gee.New()
	r.GET("/", indexHandler)
	// 设置自定义模版函数
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 注册所有模版，并注册FuncMap
	r.LoadHTMLGlob("./templates/*")

	r.Static("/assets", "./static") // 将网页中的assets/Path 映射到 ./static/Path

	r.GET("/cs", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", gee.H{
			"Now":  time.Date(2024, 12, 22, 0, 0, 0, 0, time.UTC),
			"Name": c.QueryWithDefault("name", "nil"),
		})
	})

	log.Fatal(r.Run(":8888"))
}

func test_recover() {
	r := gee.New()
	r.UseDefault()
	r.GET("/", func(c *gee.Context) {
		names := []string{"foo", "bar"}
		c.String(200, names[5])
	})

	r.GET("/h", func(c *gee.Context) {
		c.String(http.StatusOK, "ok!cc")
	})

	r.Run(":8888")
}

func main() {
	test_recover()
	fmt.Println("exec ok!")
}
