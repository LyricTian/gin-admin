package middleware

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
	"time"

	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
)

// LoggerMiddleware GIN的日志中间件
func LoggerMiddleware(allowPrefixes []string, skipPrefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if !util.CheckPrefix(p, allowPrefixes...) ||
			util.CheckPrefix(p, skipPrefixes...) {
			c.Next()
			return
		}

		nctx := context.New(c)
		method := c.Request.Method
		span := logger.StartSpan(nctx.CContext(), context.MRouterTitle[context.GetRouterTitleKey(method, p)], p)
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
				if err == nil {
					c.Request.Body.Close()
					buf := bytes.NewBuffer(body)
					c.Request.Body = ioutil.NopCloser(buf)
					fields["content_length"] = c.Request.ContentLength
					fields["body"] = string(body)
				}
			}
		}
		c.Next()

		timeConsuming := time.Since(start).Nanoseconds() / 1e6
		fields["time_consuming"] = timeConsuming
		fields["res_status"] = c.Writer.Status()
		fields["res_length"] = c.Writer.Size()
		if v, ok := c.Get(util.ContextKeyResBody); ok {
			if b, ok := v.([]byte); ok {
				fields["res_body"] = string(b)
			}
		}
		fields[logger.UserIDKey] = nctx.GetUserID()
		span.WithFields(fields).Infof("[http] %s-%s-%s-%d(%dms)",
			p, c.Request.Method, c.ClientIP(), c.Writer.Status(), timeConsuming)
	}
}
