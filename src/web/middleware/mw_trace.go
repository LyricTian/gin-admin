package middleware

import (
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
)

// TraceMiddleware 跟踪ID中间件
func TraceMiddleware(skipper SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipper != nil && skipper(c) {
			c.Next()
			return
		}

		// 优先从请求头中获取请求ID，如果没有则使用UUID
		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = util.MustUUID()
		}
		c.Set(context.ContextKeyTraceID, traceID)
		c.Next()
	}
}
