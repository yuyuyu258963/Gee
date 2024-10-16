package main

import (
	"fmt"
	gee "gee/Gee"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "index path: %v", r.URL.Path)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello\n")
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header [%s]: %v\n", k, v)
	}
}

func main() {
	r := gee.New()
	r.GET("/", indexHandler)
	r.GET("/hello", helloHandler)
	log.Fatal(r.Run(":8080"))
}
