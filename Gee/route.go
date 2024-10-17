package gee

type roteTreeRoot struct {
	root *routeTreeNode
}

// 路由树的节点，将path按照"/"进行分割，然后依次添加到路由树上
// 问题是如何区分GET和POST这些请求对应的操作，还是用map?
type routeTreeNode struct {
	nodes         map[string]*routeTreeNode
	handlerFuncs  map[string][]HandlerFunc
	middleHandler []HandlerFunc
	isLeaf        bool
}

func (r *roteTreeRoot) Init() {
	r.root = &routeTreeNode{}
	r.root.Init()
}

// 找到路径上的路由
// 如果没有找到注册的这个路由就返回
func (r *roteTreeRoot) walk(url string) *routeTreeNode {
	urlItems := splitStr(url, urlSep)
	var root *routeTreeNode = r.root
	var ok bool
	for i := range urlItems {
		if root, ok = root.nodes[urlItems[i]]; !ok {
			return nil
		}
	}
	return root
}

// 路由到最后的节点，并返回最后走到的路由节点
func (r *roteTreeRoot) walkWithCreate(url string) *routeTreeNode {
	urlItems := splitStr(url, "/")
	var root *routeTreeNode = r.root
	var tempRoot *routeTreeNode
	var ok bool
	// fmt.Println(len(urlItems), root)
	for i := range urlItems {
		if tempRoot, ok = root.nodes[urlItems[i]]; ok {
			root = tempRoot
		} else {
			newNode := &routeTreeNode{}
			newNode.Init()
			root.nodes[urlItems[i]] = newNode
			root = newNode
		}
	}
	root.isLeaf = true
	return root
}

func (r *routeTreeNode) Init() {
	r.nodes = make(map[string]*routeTreeNode)
	r.handlerFuncs = make(map[string][]HandlerFunc, 0)
	r.middleHandler = make([]HandlerFunc, 0)
}

func (r *routeTreeNode) addHandler(method string, handler HandlerFunc) {
	if handlers, ok := r.handlerFuncs[method]; ok {
		handlers = append(handlers, handler)
		r.handlerFuncs[method] = handlers
	} else {
		r.handlerFuncs[method] = []HandlerFunc{handler}
	}
}

func (r *routeTreeNode) addMiddleHandler(handler HandlerFunc) {
	r.middleHandler = append(r.middleHandler, handler)
}

// 调用路由节点上的所有handler
func (r *routeTreeNode) callHandler(method string, c *Context) {
	for i := range r.handlerFuncs[method] {
		r.handlerFuncs[method][i](c)
	}
}
