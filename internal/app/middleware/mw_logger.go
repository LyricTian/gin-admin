package middleware

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/ginplus"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		method := c.Request.Method
		span := logger.StartSpan(ginplus.NewContext(c),
			logger.SetSpanTitle("访问日志"),
			logger.SetSpanFuncName(JoinRouter(method, p)))

		start := time.Now()

		fields := make(map[string]interface{})
		fields["ip"] = c.ClientIP()
		fields["method"] = method
		fields["url"] = c.Request.URL.String()
		fields["proto"] = c.Request.Proto
		fields["header"] = c.Request.Header
		fields["user_agent"] = c.GetHeader("User-Agent")

		// 如果是POST/PUT请求，并且内容类型为JSON，则读取内容体
		if method == http.MethodPost || method == http.MethodPut {
			mediaType, _, _ := mime.ParseMediaType(c.GetHeader("Content-Type"))
			if mediaType == "application/json" {
				body, err := ioutil.ReadAll(c.Request.Body)
				c.Request.Body.Close()
				if err == nil {
					buf := bytes.NewBuffer(body)
					c.Request.Body = ioutil.NopCloser(buf)
					fields["content_length"] = c.Request.ContentLength
					fields["body"] = string(body)
				}
			}
		}
		c.Next()

		timeConsuming := time.Since(start).Nanoseconds() / 1e6
		fields["res_status"] = c.Writer.Status()
		fields["res_length"] = c.Writer.Size()

		if v, ok := c.Get(ginplus.ResBodyKey); ok {
			if b, ok := v.([]byte); ok {
				fields["res_body"] = string(b)
			}
		}

		fields[logger.UserIDKey] = ginplus.GetUserID(c)
		span.WithFields(fields).Infof("[http] %s-%s-%s-%d(%dms)",
			p, c.Request.Method, c.ClientIP(), c.Writer.Status(), timeConsuming)
	}
}
