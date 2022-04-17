package hee

import (
	"net/http"
)

//类型适配器
// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix string // 当前组的公共前缀
	//middleware []HandlerFunc // 支持中间件
	parent *RouterGroup //支持分组嵌套
	engine *Engine      //共享同一个Engine
}

type Engine struct {
	*RouterGroup                //最顶层的路由分组
	router       *router        //储存所有的路由
	groups       []*RouterGroup //保存所有的路由分组
}

//初始化一个hee的实例
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//创建一个RouterGroup
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	//生成一个新的group
	newRouterGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newRouterGroup)
	return newRouterGroup
}

//往路由器中添加路由
func (g *RouterGroup) addRoute(method string, pattern string, handle HandlerFunc) {
	pattern = g.prefix + pattern
	g.engine.router.addRoute(method, pattern, handle)
}

//GET请求
func (g *RouterGroup) Get(pattern string, handle HandlerFunc) {
	g.addRoute("GET", pattern, handle)
}

//POST请求
func (g *RouterGroup) POST(pattern string, handle HandlerFunc) {
	g.addRoute("POST", pattern, handle)
}

//还有其它Restful风格方法。。。。

//使Engine实现http.Handler接口
//在这里实现上下文封装
func (g *RouterGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	g.engine.router.handle(c)
}

func (g *RouterGroup) Run(addr string) (err error) {
	return http.ListenAndServe(addr, g)
}
