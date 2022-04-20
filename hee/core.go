package hee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

//类型适配器
// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string        // 当前组的公共前缀
	middlewares []HandlerFunc // 支持中间件
	parent      *RouterGroup  //支持分组嵌套
	engine      *Engine       //共享同一个Engine
}

type Engine struct {
	*RouterGroup                     //最顶层的路由分组
	router        *router            //储存所有的路由
	groups        []*RouterGroup     //保存所有的路由分组
	htmlTemplates *template.Template // 将所有的模板加载进内存
	funcMap       template.FuncMap   // 所有的自定义模板渲染函数
}

//初始化一个hee的实例
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
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

//将中间件应用到group
func (g *RouterGroup) Use(middleware ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}

//往路由器中添加路由
func (g *RouterGroup) addRoute(method string, pattern string, handle HandlerFunc) {
	g.engine.router.addRoute(method, g.prefix+pattern, handle)
}

//GET请求
func (g *RouterGroup) GET(pattern string, handle HandlerFunc) {
	g.addRoute("GET", pattern, handle)
}

//POST请求
func (g *RouterGroup) POST(pattern string, handle HandlerFunc) {
	g.addRoute("POST", pattern, handle)
}

//还有其它Restful风格方法。。。。

//创建静态文件处理handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absloutePath := path.Join(group.prefix, relativePath)
	//去除absloutePath
	fileServer := http.StripPrefix(absloutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.GetParam("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.SetStatusCode(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

//解析请求的地址，映射到服务器上文件的真实地址
func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	g.GET(urlPattern, handler)
}

//使Engine实现http.Handler接口
//在这里实现上下文封装
//请求处理的入口
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	//对所有的group进行匹配
	for _, group := range e.groups {
		//请求path在路由器中有指定前缀，将该前缀对应的middleware加入到该请求中
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.engine = e // 初始化engine
	c.handlers = middlewares
	e.router.handle(c)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
