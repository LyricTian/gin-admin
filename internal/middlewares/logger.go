package middlewares

import (
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoggerConfig struct {
	AllowedPathPrefixes      []string
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
		if !AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			SkippedPathPrefixes(c, config.SkippedPathPrefixes...) {
			c.Next()
			return
		}

		start := time.Now()
		contentType := c.Request.Header.Get("Content-Type")

		fields := []zap.Field{
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
			zap.String("uri", c.Request.RequestURI),
			zap.String("host", c.Request.Host),
			zap.String("remote_addr", c.Request.RemoteAddr),
			zap.String("proto", c.Request.Proto),
			zap.Int64("content_length", c.Request.ContentLength),
			zap.String("content_type", contentType),
			zap.String("pragma", c.Request.Header.Get("Pragma")),
		}

		c.Next()

		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			mediaType, _, _ := mime.ParseMediaType(contentType)
			if mediaType == "application/json" {
				if v, ok := c.Get(utils.RequestBodyKey); ok {
					if b, ok := v.([]byte); ok && len(b) <= config.MaxOutputRequestBodyLen {
						fields = append(fields, zap.String("body", string(b)))
					}
				}
			}
		}

		cost := time.Since(start).Nanoseconds() / 1e6
		fields = append(fields, zap.Int64("cost", cost))
		fields = append(fields, zap.Int("status", c.Writer.Status()))
		fields = append(fields, zap.String("res_time", time.Now().Format("2006-01-02 15:04:05.999")))
		fields = append(fields, zap.Int("res_size", c.Writer.Size()))

		if v, ok := c.Get(utils.ResponseBodyKey); ok {
			if b, ok := v.([]byte); ok && len(b) <= config.MaxOutputResponseBodyLen {
				fields = append(fields, zap.String("res_body", string(b)))
			}
		}

		ctx := c.Request.Context()
		ctx = logging.NewTag(ctx, logging.TagKeyRequest)
		logging.Context(ctx).Info(fmt.Sprintf("[HTTP] %s-%s-%d (%dms)",
			c.Request.URL.Path, c.Request.Method, c.Writer.Status(), cost), fields...)
	}
}
