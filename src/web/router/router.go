package router

import (
	"path"

	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
)

// MRouterTitle 路由关联的标题数据
var MRouterTitle = make(map[string]string)

// HandlerFunc 处理函数
type HandlerFunc func(*context.Context)

// Handle registers a new request handle and middleware with the given path and method.
func Handle(g *gin.RouterGroup, httpMethod string, relativePath string, handler HandlerFunc, title string) {
	titleKey := path.Join(httpMethod, g.BasePath(), relativePath)
	MRouterTitle[titleKey] = title
	g.Handle(httpMethod, relativePath, func(c *gin.Context) {
		handler(context.NewContext(c))
	})
}

// GET is a shortcut for router.Handle("GET", path, handle).
func GET(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "GET", relativePath, handler, title)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func POST(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "POST", relativePath, handler, title)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func DELETE(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "DELETE", relativePath, handler, title)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func PATCH(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "PATCH", relativePath, handler, title)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func PUT(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "PUT", relativePath, handler, title)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func OPTIONS(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "OPTIONS", relativePath, handler, title)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func HEAD(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "HEAD", relativePath, handler, title)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func Any(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "GET", relativePath, handler, title)
	Handle(g, "POST", relativePath, handler, title)
	Handle(g, "PUT", relativePath, handler, title)
	Handle(g, "PATCH", relativePath, handler, title)
	Handle(g, "HEAD", relativePath, handler, title)
	Handle(g, "OPTIONS", relativePath, handler, title)
	Handle(g, "DELETE", relativePath, handler, title)
	Handle(g, "CONNECT", relativePath, handler, title)
	Handle(g, "TRACE", relativePath, handler, title)
}
