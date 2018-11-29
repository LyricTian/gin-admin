package router

import (
	"strings"

	"github.com/LyricTian/gin-admin/src/util"
	"github.com/gin-gonic/gin"
)

// TraceMiddleware 跟踪ID中间件
func TraceMiddleware(allowPrefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		allow := false
		for _, p := range allowPrefixes {
			if strings.HasPrefix(c.Request.URL.Path, p) {
				allow = true
				break
			}
		}

		if !allow {
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
