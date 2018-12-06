package router

import (
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/gin-gonic/gin"
)

// TraceMiddleware 跟踪ID中间件
func TraceMiddleware(allowPrefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !util.CheckPrefix(c.Request.URL.Path, allowPrefixes...) {
			c.Next()
			return
		}

		traceID := c.Query("X-Request-Id")
		if traceID == "" {
			traceID = util.MustUUID()
		}
		c.Set(util.ContextKeyTraceID, traceID)
		c.Next()
	}
}
