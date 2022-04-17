package hee

import (
	"net/http"
)

//类型适配器
// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type Engine struct {
	router *router //储存所有的路由
}

//初始化一个hee的实例
func New() *Engine {
	return &Engine{router: newRouter()}
}

//往路由器中添加路由
func (e *Engine) addRoute(method string, pattern string, handle HandlerFunc) {
	e.router.addRoute(method, pattern, handle)
}

//GET请求
func (e *Engine) Get(pattern string, handle HandlerFunc) {
	e.addRoute("GET", pattern, handle)
}

//POST请求
func (e *Engine) POST(pattern string, handle HandlerFunc) {
	e.addRoute("POST", pattern, handle)
}

//还有其它Restful风格方法。。。。

//使Engine实现http.Handler接口
//在这里实现上下文封装
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
