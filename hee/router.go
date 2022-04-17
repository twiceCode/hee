package hee

import (
	"net/http"
	"strings"
)

type router struct {
	root     map[string]*tireNode
	handlers map[string]HandlerFunc
}

//初始化一个路由器
func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		root:     make(map[string]*tireNode),
	}
}

//解析路由（只能处理一个*），划分为一段段
func parsePattern(pattern string) (parts []string) {
	pt := strings.Split(pattern, "/")

	for _, v := range pt {
		if v != "" {
			parts = append(parts, v)
			//如果匹配到了某一段的第一个字符位*,后续不必再进行
			if v[0] == '*' {
				break
			}
		}
	}
	return
}

//添加一个路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {

	parts := parsePattern(pattern)
	key := method + "-" + pattern
	//查询
	if _, ok := r.root[method]; !ok {
		r.root[method] = &tireNode{}
	}
	r.root[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

//获取匹配路由
func (r *router) getRoute(method string, path string) (*tireNode, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	node, ok := r.root[method]
	if !ok {
		return nil, nil
	}
	n := node.search(searchParts, 0)
	if n != nil {
		//获取路由params参数
		parts := parsePattern(n.pattern)
		for i, v := range parts {
			if v[0] == ':' {
				params[v[1:]] = searchParts[i]
			}
			if v[0] == '*' && len(v) > 1 {
				params[v[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		key := c.Method + "-" + node.pattern
		c.Params = params
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
