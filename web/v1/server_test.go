package v1

import (
	"net/http"
	"testing"
)

// // 方案一
// func TestServer(t *testing.T) {
// 	var s Server
// 	// Server 组合了http.Handler，即实现了http.Handler接口
// 	// 故直接传变量`s`即可
// 	http.ListenAndServe(":8080", s)
// }

// 方案二
func TestServer(t *testing.T) {
	var s Server

	// 当作普通的http.Handler使用
	http.ListenAndServe(":8080", s)
	http.ListenAndServeTLS(":8080", "create file", "key file", s)

	// 当作独立的实体使用
	s.Start(":8081")
}
