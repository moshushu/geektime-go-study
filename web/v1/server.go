package v1

import (
	"net"
	"net/http"
)

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
	http.Handler
	// 可以在start方法中，利用net.Listen()监听端口，使用http.Serve()来启动服务
	// 同时也可在start做defer start的回调，已实现能自主管理生命周期能力
	Start(add string) error
}

// 检测HTTPServer是否实现Server接口
var _ Server = &HTTPServer{}

// HTTPServer 是为了实现Server接口的实例，并依托Server接口增添新功能
type HTTPServer struct{}

// ServeHTTP 整个Web核心入口
// ServeHTTP 为了实现http.Handler接口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {}

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
