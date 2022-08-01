package middleware

import (
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoggerConfig struct {
	SkippedPathPrefixes      []string
	MaxOutputRequestBodyLen  int
	MaxOutputResponseBodyLen int
}

var DefaultLoggerConfig = LoggerConfig{
	MaxOutputRequestBodyLen:  1024 * 1024,
	MaxOutputResponseBodyLen: 1024 * 1024,
}

// Record detailed request logs for quick troubleshooting.
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

func LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkippedPathPrefixes(c, config.SkippedPathPrefixes...) {
			c.Next()
			return
		}

		start := time.Now()
		p := c.Request.URL.Path
		method := c.Request.Method
		contentType := c.Request.Header.Get("Content-Type")

		var fields []zap.Field
		fields = append(fields, zap.String("client_ip", c.ClientIP()))
		fields = append(fields, zap.String("method", method))
		fields = append(fields, zap.String("path", p))
		fields = append(fields, zap.String("user_agent", c.Request.UserAgent()))
		fields = append(fields, zap.String("referer", c.Request.Referer()))
		fields = append(fields, zap.String("uri", c.Request.RequestURI))
		fields = append(fields, zap.String("host", c.Request.Host))
		fields = append(fields, zap.String("remote_addr", c.Request.RemoteAddr))
		fields = append(fields, zap.String("proto", c.Request.Proto))
		fields = append(fields, zap.Int64("content_length", c.Request.ContentLength))
		fields = append(fields, zap.String("content_type", contentType))
		fields = append(fields, zap.String("pragma", c.Request.Header.Get("Pragma")))

		if method == http.MethodPost || method == http.MethodPut {
			mediaType, _, _ := mime.ParseMediaType(contentType)
			if mediaType != "multipart/form-data" {
				if v, ok := c.Get(utilx.RequestBodyKey); ok {
					if b, ok := v.([]byte); ok && len(b) <= config.MaxOutputRequestBodyLen {
						fields = append(fields, zap.String("body", string(b)))
					}
				}
			}
		}

		c.Next()

		cost := time.Since(start).Nanoseconds() / 1e6
		fields = append(fields, zap.Int64("cost", cost))
		fields = append(fields, zap.Int("status", c.Writer.Status()))
		fields = append(fields, zap.String("res_time", time.Now().Format("2006-01-02 15:04:05.999")))
		fields = append(fields, zap.Int("res_size", c.Writer.Size()))

		if v, ok := c.Get(utilx.ResponseBodyKey); ok {
			if b, ok := v.([]byte); ok && len(b) <= config.MaxOutputResponseBodyLen {
				fields = append(fields, zap.String("res_body", string(b)))
			}
		}

		ctx := c.Request.Context()
		ctx = logger.NewTag(ctx, logger.TagKeyRequest)
		logger.Context(ctx).Info(fmt.Sprintf("[http] %s-%s-%d (%dms)", p, method, c.Writer.Status(), cost), fields...)
	}
}
