package router

import (
	"gin-admin/src/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// TraceMiddleware 跟踪ID中间件
func TraceMiddleware(prefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		allow := false
		for _, p := range prefixes {
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
			traceID = util.UUIDString()
		}
		c.Set(util.ContextKeyTraceID, traceID)
		c.Next()
	}
}
