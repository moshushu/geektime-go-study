package v1

import "net/http"

// Context 需要*http.Request 和 http.ResponseWriter 是由于ServeHTTP定的
type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
}
