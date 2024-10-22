package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc Used by Gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHttp
type Engine struct {
	*RouteGroup   // index RouteGroup
	route         *router
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap
}

// 因为routeGroup中的操作和Engine上的方法有很多重叠所以写到这个文件中
type RouteGroup struct {
	prefix      string //包含从到当前路由的前缀
	parent      *RouteGroup
	nextGroups  []*RouteGroup
	middlewares []HandlerFunc
	engine      *Engine //因为后序需要直接通过group对象进行操作，但是还需要原始的框架的能力
}

func newEngine() *Engine {
	e := &Engine{route: newRouter()}
	e.RouteGroup = newRouteGroup("", nil, e) //总控的前缀为空字符串
	return e
}

// 创建一个路由分组
func newRouteGroup(prefix string, parent *RouteGroup, e *Engine) *RouteGroup {
	return &RouteGroup{
		prefix:      prefix,
		nextGroups:  make([]*RouteGroup, 0),
		middlewares: make([]HandlerFunc, 0),
		parent:      parent,
		engine:      e,
	}
}

// TODO:
// implement ListenAdnServe interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(req, w) // 每个请求处理的开始时创建一个上下文
	c.engine = e

	middlewares := e.collectMiddlewares(c.Path)
	c.handles = append(c.handles, middlewares...)
	fn := e.route.getHandle(c)
	c.handles = append(c.handles, fn) // 表示中间件执行结束后再执行路由树上查找到的处理函数
	c.Next()
}

// set template func
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

// 将所有的模版都添加到内存中
func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

// run and listen request at a port
func (e *Engine) Run(port string) error {
	return http.ListenAndServe(port, e)
}

func (rg *RouteGroup) GET(path string, handler HandlerFunc) {
	rg.engine.route.addRoute("GET", rg.prefix+path, handler)
}

func (rg *RouteGroup) POST(path string, handler HandlerFunc) {
	rg.engine.route.addRoute("POST", rg.prefix+path, handler)
}

// 嵌套地添加分组
func (rg *RouteGroup) Group(prefix string) (g *RouteGroup) {
	if rg.engine.RouteGroup != rg {
		prefix = rg.prefix + "/" + prefix
	}

	if prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	g = newRouteGroup(prefix, rg, rg.engine)
	rg.nextGroups = append(rg.nextGroups, g)
	return g
}

// 在分组控件上新增中间件
func (rg *RouteGroup) Use(fn ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, fn...)
}

// 根据请求路径收集所有的中间件
func (rg *RouteGroup) collectMiddlewares(pattern string) []HandlerFunc {
	ws := rg.middlewares
	for _, w := range rg.nextGroups {
		if strings.HasPrefix(pattern, w.prefix) { /// 递归地获得所有的中间件
			ws = append(ws, w.collectMiddlewares(pattern)...)
		}
	}
	return ws
}

// 静态资源处理逻辑
func (rg *RouteGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(rg.prefix, relativePath)
	// 类似于http.ListenAndServe， http.StripPrefix就是给这个请求前缀的url通过传入的handler处理
	// 所以这个 absolutePath 指的是网页上请求fs文件系统使用的路径是绝对的
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// check if the file exists and if we permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req) // 返回结果
	}
}

// serve static files
// 将磁盘上的某个文件夹root映射到relativePath
// r := gee.New()
// r.Static("/assets", "./static")
// r.Run(":9999")
// 用户访问 "/assets/hom/doc" 就对应 "./static/hom/doc"
func (rg *RouteGroup) Static(relativePath string, root string) {
	handler := rg.createStaticHandler(relativePath, http.Dir(root))
	// 注册路由
	urlPattern := path.Join(relativePath, "/*filepath")
	rg.GET(urlPattern, handler)
}

// New is the constructor of gee.Engine
func New() *Engine {
	return newEngine()
}
