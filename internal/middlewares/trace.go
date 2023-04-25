package middlewares

import (
	"fmt"
	"strings"

	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type TraceConfig struct {
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
	RequestHeaderKey    string
	ResponseTraceKey    string
}

var DefaultTraceConfig = TraceConfig{
	RequestHeaderKey: "X-Request-Id",
	ResponseTraceKey: "X-Trace-Id",
}

func Trace() gin.HandlerFunc {
	return TraceWithConfig(DefaultTraceConfig)
}

func TraceWithConfig(config TraceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			SkippedPathPrefixes(c, config.SkippedPathPrefixes...) {
			c.Next()
			return
		}

		traceID := c.GetHeader(config.RequestHeaderKey)
		if traceID == "" {
			traceID = fmt.Sprintf("TRACE-%s", strings.ToUpper(xid.New().String()))
		}

		ctx := utils.NewTraceID(c.Request.Context(), traceID)
		ctx = logging.NewTraceID(ctx, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(config.ResponseTraceKey, traceID)
		c.Next()
	}
}
