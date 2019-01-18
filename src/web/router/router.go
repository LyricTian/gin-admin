package router

import (
	"path"

	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
)

var defaultOptions = options{}

type options struct {
	title string
}

// Option 定义配置项
type Option func(*options)

// SetTitle 设定路由标题
func SetTitle(title string) Option {
	return func(o *options) {
		o.title = title
	}
}

// HandlerFunc 处理函数
type HandlerFunc func(*context.Context)

// Handle registers a new request handle and middleware with the given path and method.
func Handle(g *gin.RouterGroup, httpMethod string, relativePath string, handler HandlerFunc, opts ...Option) {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	context.SetRouterTitle(httpMethod, path.Join(g.BasePath(), relativePath), o.title)
	g.Handle(httpMethod, relativePath, func(c *gin.Context) {
		handler(context.New(c))
	})
}

// GET is a shortcut for router.Handle("GET", path, handle).
func GET(g *gin.RouterGroup, relativePath string, handler HandlerFunc, opts ...Option) {
	Handle(g, "GET", relativePath, handler, opts...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func POST(g *gin.RouterGroup, relativePath string, handler HandlerFunc, opts ...Option) {
	Handle(g, "POST", relativePath, handler, opts...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func DELETE(g *gin.RouterGroup, relativePath string, handler HandlerFunc, opts ...Option) {
	Handle(g, "DELETE", relativePath, handler, opts...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func PATCH(g *gin.RouterGroup, relativePath string, handler HandlerFunc, opts ...Option) {
	Handle(g, "PATCH", relativePath, handler, opts...)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func PUT(g *gin.RouterGroup, relativePath string, handler HandlerFunc, opts ...Option) {
	Handle(g, "PUT", relativePath, handler, opts...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func OPTIONS(g *gin.RouterGroup, relativePath string, handler HandlerFunc, opts ...Option) {
	Handle(g, "OPTIONS", relativePath, handler, opts...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func HEAD(g *gin.RouterGroup, relativePath string, handler HandlerFunc, opts ...Option) {
	Handle(g, "HEAD", relativePath, handler, opts...)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func Any(g *gin.RouterGroup, relativePath string, handler HandlerFunc, opts ...Option) {
	Handle(g, "GET", relativePath, handler, opts...)
	Handle(g, "POST", relativePath, handler, opts...)
	Handle(g, "PUT", relativePath, handler, opts...)
	Handle(g, "PATCH", relativePath, handler, opts...)
	Handle(g, "HEAD", relativePath, handler, opts...)
	Handle(g, "OPTIONS", relativePath, handler, opts...)
	Handle(g, "DELETE", relativePath, handler, opts...)
	Handle(g, "CONNECT", relativePath, handler, opts...)
	Handle(g, "TRACE", relativePath, handler, opts...)
}
