package gee

type roteTreeRoot struct {
	root          *routeTreeNode
	handlerFunc   []HandlerFunc
	middleHandler []HandlerFunc
}

// 路由树的节点，将path按照"/"进行分割，然后依次添加到路由树上
type routeTreeNode struct {
	next          map[string]*routeTreeNode
	handlerFunc   []HandlerFunc
	middleHandler []HandlerFunc
}

var route roteTreeRoot
