package middleware

import (
	"mime"
	"net/http"
	"time"

	"github.com/LyricTian/gin-admin/v8/internal/app/config"
	"github.com/LyricTian/gin-admin/v8/internal/app/ginx"
	"github.com/LyricTian/gin-admin/v8/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Request logger
func LoggerMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		method := c.Request.Method

		start := time.Now()
		fields := make(map[string]interface{})
		fields["ip"] = c.ClientIP()
		fields["method"] = method
		fields["url"] = c.Request.URL.String()
		fields["proto"] = c.Request.Proto
		fields["header"] = c.Request.Header
		fields["user_agent"] = c.GetHeader("User-Agent")
		fields["content_length"] = c.Request.ContentLength

		if method == http.MethodPost || method == http.MethodPut {
			mediaType, _, _ := mime.ParseMediaType(c.GetHeader("Content-Type"))
			if mediaType != "multipart/form-data" {
				if v, ok := c.Get(ginx.ReqBodyKey); ok {
					if b, ok := v.([]byte); ok && len(b) <= config.C.HTTP.MaxReqLoggerLength {
						fields["body"] = string(b)
					}
				}
			}
		}
		c.Next()

		timeConsuming := time.Since(start).Nanoseconds() / 1e6
		fields["res_status"] = c.Writer.Status()
		fields["res_length"] = c.Writer.Size()

		if v, ok := c.Get(ginx.ResBodyKey); ok {
			if b, ok := v.([]byte); ok && len(b) <= config.C.HTTP.MaxResLoggerLength {
				fields["res_body"] = string(b)
			}
		}

		ctx := c.Request.Context()
		entry := logger.WithContext(logger.NewTagContext(ctx, "__request__"))
		entry.WithFields(fields).Infof("[http] %s-%s-%s-%d(%dms)",
			p, c.Request.Method, c.ClientIP(), c.Writer.Status(), timeConsuming)
	}
}
