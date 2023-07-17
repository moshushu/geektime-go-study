package v1

import (
	"net"
	"net/http"
)

// HandleFunc 处理逻辑
type HandleFunc func(ctx *Context)

// Server 接口，服务器抽象
// 从Server特性上来说，至少提供三部分功能：
// 1、生命周期控制：启动、关闭等
// 2、路由注册接口：提供路由注册功能
// 3、作为http包到Web框架的桥梁，所有的Web框架都是基于http包来进行的；即实现http.Handler接口

// 方案一：
// type Server interface {
// 	// 只组合http.Handler
// 	// 与原生的http包使用无异，只需调用http.ListenAndServe即可，但很难实现生命周期的控制
// 	http.Handler
// }

// 方案二
type Server interface {
	// 组合http.Handler，并且增加Start方法作为服务启动
	// 可以当作普通的http.Handler使用，又可以作为一个独立实体使用
	// 拥有自己的管理生命周期的能力
	http.Handler // 1、作为http包到Web框架的桥梁（三部分功能之一）
	// 可以在start方法中，利用net.Listen()监听端口，使用http.Serve()来启动服务
	// 同时也可在start做defer start的回调，已实现能自主管理生命周期能力
	Start(add string) error // 2、生命周期控制（三部分功能之二）

	// 路由注册
	// 路由注册有两种方式：1、只允许注册一个HandleFunc；2、允许注册不定个HandleFunc
	// 只允许注册一个（推荐）
	AddRoute(method string, path string, handler HandleFunc)
	// 允许注册不定个HandleFunc
	AddRoutes(method string, path string, handler ...HandleFunc)
	// 解释：为什么推荐只允许注册一个？
	// 允许注册多个HandleFunc带来的考虑是，我什么时候要停止下来？其中一个失败了，是否继续执行下去？如果需要中断某个HandleFunc应该怎么中断？
	// 这些都是需要格外的复杂处理，但这部分是可以交由用户来处理的，我们可以只允许注册一个，即便真有多个的场景，用户可以自己组合成一个
}

// 检测HTTPServer是否实现Server接口
var _ Server = &HTTPServer{}

// HTTPServer 是为了实现Server接口的实例，并依托Server接口增添新功能
type HTTPServer struct{}

// ServeHTTP 整个Web核心入口
// ServeHTTP 为了实现http.Handler接口
// 包含的功能：Context构建、路由匹配、执行业务逻辑
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Context构建
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 执行业务逻辑
	h.serve(ctx)
}
func (h *HTTPServer) serve(*Context) {}

func (h *HTTPServer) Start(add string) error {
	l, err := net.Listen("tcp", add)
	if err != nil {
		return err
	}

	// 在这里，可以让用户注册所谓的 after start 回调
	// 比如说往你的 admin 注册一下自己这个实例
	// 在这里执行一些你业务所需的前置条件
	// 所谓的生命周期就可以在这实现

	return http.Serve(l, h)
}

func (h *HTTPServer) AddRoutes(method string, path string, handler ...HandleFunc) {}

func (h *HTTPServer) AddRoute(method string, path string, handler HandleFunc) {}

// 可以基于AddRoute提供更多的方法，如GET、POST等等
// 为什么不这些方法放到接口里面，而是实现（实例）里面？主要原因是保持核心接口Server的少而美
// 这部分属于是衍生api，而AddRoute 属于是核心api
func (h *HTTPServer) Get(path string, handler HandleFunc) {
	h.AddRoute(http.MethodGet, path, handler)
}
func (h *HTTPServer) Post(path string, handler HandleFunc) {
	h.AddRoute(http.MethodPost, path, handler)
}
