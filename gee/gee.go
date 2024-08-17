package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// 路由映射处理方法
type HandleFunc func(*Context)

// 映射路由表router, key为请求方法+请求路径，value为处理方法
type Engine struct {
	router        *router
	*RouterGroup                     // 继承RouterGroup的属性
	groups        []*RouterGroup     // 存储所有的group
	htmlTemplates *template.Template // 解析模板
	funcMap       template.FuncMap   // 自定义模板函数
}

// 路由组,实现分组添加路由
type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc // 支持中间件
	engine      *Engine      // 所有的group共享一个engine实例
	parent      *RouterGroup // 支持嵌套
}

// 创建默认引擎,并添加默认中间件
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 创建路由组
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRoute(method string, pattern string, handler HandleFunc) {
	pattern = g.prefix + pattern
	log.Printf("Route %4s - %s", method, pattern)
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(pattern string, handler HandleFunc) {
	g.addRoute("GET", pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandleFunc) {
	g.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 添加中间件
func (g *RouterGroup) Use(middlewares ...HandleFunc) {
	g.middlewares = append(g.middlewares, middlewares...)

}

// 解析请求，查找路由映射表
// 通过URL前缀查找对应的中间件
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandleFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...) // 收集所有的中间件
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

// 创建静态文件服务
func (g *RouterGroup) CreateStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// 检查文件是否存在
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.writer, c.Req)
	}
}

// 添加静态文件服务
func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.CreateStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	g.GET(urlPattern, handler)
}

// 设置模板
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

// 加载模板
func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}
