package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"time"

	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

		start := time.Now()
		fields := logrus.Fields{}
		fields["ip"] = c.ClientIP()
		fields["method"] = c.Request.Method
		fields["url"] = c.Request.URL.String()
		fields["proto"] = c.Request.Proto
		fields["header"] = c.Request.Header
		fields["user_agent"] = c.GetHeader("User-Agent")

		if m := c.Request.Method; m == http.MethodPost ||
			m == http.MethodPut {
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

		fields["time"] = fmt.Sprintf("%dms", time.Since(start).Nanoseconds()/1e6)
		fields["status"] = c.Writer.Status()
		fields["length"] = c.Writer.Size()

		m := p
		if v := c.GetString(util.ContextKeyURLMemo); v != "" {
			m = fmt.Sprintf("%s(%s)", p, v)
		}

		logger.Access(
			c.GetString(util.ContextKeyTraceID),
			c.GetString(util.ContextKeyUserID),
		).WithFields(fields).Infof(
			"[http] %s - %s - %s",
			m,
			c.Request.Method,
			c.ClientIP(),
		)
	}
}
