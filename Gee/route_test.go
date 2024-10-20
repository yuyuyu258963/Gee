package gee

import (
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/c", nil)
	r.addRoute("GET", "/home/:name", nil)
	r.addRoute("GET", "/home/b/c", nil)
	r.addRoute("GET", "/a/*filepath", nil)
	return r
}
func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("p/:name"), []string{"p", ":name"})
	ok = reflect.DeepEqual(parsePattern("p/*"), []string{"p", "*"})
	ok = reflect.DeepEqual(parsePattern("p/*name/a"), []string{"p", "*name"})
	ok = reflect.DeepEqual(parsePattern("p/a"), []string{"p", "a"})
	if !ok {
		t.Fatal("test parsePatten Error")
	}
}

func TestRoute(t *testing.T) {
	r := newTestRouter()
	n, params := r.getRoute("GET", "/home/ywh")
	if n == nil {
		t.Fatalf("route cant found ")
	}
	if n.pattern != "/home/:name" {
		t.Fatal("Error route pattern")
	}

	if params["name"] != "ywh" {
		t.Fatal("can't get aim param")
	}

	n, _ = r.getRoute("GET", "/a/b/c")
	if n == nil {
		t.Fatal("can not found a route path")
	}
	if n.pattern != "/a/*filepath" {
		t.Fatal("can't get aim pattern'")
	}
}
