package logger

import (
	"bytes"
	"gin-admin/src/util"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 定义日志中使用的键名
const (
	FieldKeyType   = "type"
	FieldKeyUserID = "user_id"
)

var defaultOptions = options{
	level:  5,
	format: "text",
}

var internalLogger *Logger
var once sync.Once

type options struct {
	level  int
	format string
}

// Option 定义配置参数
type Option func(o *options)

// SetLevel 设定日志级别(0:panic,1:fatal,2:error,3:warn,4:info,5:debug)
func SetLevel(level int) Option {
	return func(o *options) {
		o.level = level
	}
}

// SetFormat 设定日志格式(text/json)
func SetFormat(format string) Option {
	return func(o *options) {
		o.format = format
	}
}

func logger() *Logger {
	if internalLogger == nil {
		internalLogger = New()
	}
	return internalLogger
}

// Default 获取默认日志实例
func Default() *Logger {
	return logger()
}

// System 系统日志
func System(userID ...string) *logrus.Entry {
	return logger().System(userID...)
}

// Access 访问日志
func Access(userID ...string) *logrus.Entry {
	return logger().Access(userID...)
}

// Operate 操作日志
func Operate(userID ...string) *logrus.Entry {
	return logger().Operate(userID...)
}

// Login 登录(登出)日志
func Login(userID string) *logrus.Entry {
	return logger().Login(userID)
}

// New 创建日志实例
func New(opts ...Option) *Logger {
	once.Do(func() {
		o := defaultOptions
		for _, opt := range opts {
			opt(&o)
		}

		l := logrus.New()
		l.SetLevel(logrus.Level(o.level))
		if o.format == "json" {
			l.Formatter = new(logrus.JSONFormatter)
		}
		internalLogger = &Logger{l}
	})
	return internalLogger
}

// HookFlusher 将缓冲区数据写入日志钩子完成接口
type HookFlusher interface {
	Flush()
}

// Logger 日志管理
type Logger struct {
	*logrus.Logger
}

func (a *Logger) typeEntry(fieldType string, userID ...string) *logrus.Entry {
	fields := logrus.Fields{
		FieldKeyType: fieldType,
	}
	if len(userID) > 0 {
		fields[FieldKeyUserID] = userID[0]
	}
	return a.WithFields(fields)
}

// System 系统日志
func (a *Logger) System(userID ...string) *logrus.Entry {
	return a.typeEntry("system", userID...)
}

// Access 访问日志
func (a *Logger) Access(userID ...string) *logrus.Entry {
	return a.typeEntry("access", userID...)
}

// Operate 操作日志
func (a *Logger) Operate(userID ...string) *logrus.Entry {
	return a.typeEntry("operate", userID...)
}

// Login 登录(登出)日志
func (a *Logger) Login(userID string) *logrus.Entry {
	return a.typeEntry("login", userID)
}

// Middleware GIN的日志中间件
func Middleware(prefixes ...string) gin.HandlerFunc {
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

		start := time.Now()
		fields := logrus.Fields{}
		fields["ip"] = c.ClientIP()
		fields["method"] = c.Request.Method
		fields["url"] = c.Request.URL.String()
		fields["proto"] = c.Request.Proto
		fields["user_agent"] = c.GetHeader("User-Agent")

		if method := c.Request.Method; method == http.MethodPost || method == http.MethodPut {
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

		fields["time_consuming"] = time.Since(start) / 1e6
		fields["status"] = c.Writer.Status()
		fields["length"] = c.Writer.Size()

		logger().Access(c.GetString(util.ContextKeyUserID)).
			WithFields(fields).
			Infof("[http] %s(%s) - %s - %s",
				c.Request.URL.Path,
				c.GetString(util.ContextKeyURLMemo),
				c.Request.Method, c.ClientIP())
	}
}
