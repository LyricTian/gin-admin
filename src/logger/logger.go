package logger

import (
	"context"
	"sync"

	"github.com/LyricTian/gin-admin/src/util"
	"github.com/sirupsen/logrus"
)

// 定义日志中使用的键名
const (
	FieldKeyTraceID = "trace_id"
	FieldKeyType    = "type"
	FieldKeyUserID  = "user_id"
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
	return New()
}

// Default 获取默认日志实例
func Default() *Logger {
	return logger()
}

// SystemWithContext 系统日志
func SystemWithContext(ctx context.Context) *logrus.Entry {
	return System(util.FromTraceIDContext(ctx), util.FromUserIDContext(ctx))
}

// System 系统日志
func System(traceID string, userID ...string) *logrus.Entry {
	return logger().System(traceID, userID...)
}

// AccessWithContext 访问日志
func AccessWithContext(ctx context.Context) *logrus.Entry {
	return Access(util.FromTraceIDContext(ctx), util.FromUserIDContext(ctx))
}

// Access 访问日志
func Access(traceID string, userID ...string) *logrus.Entry {
	return logger().Access(traceID, userID...)
}

// OperateWithContext 操作日志
func OperateWithContext(ctx context.Context) *logrus.Entry {
	return Operate(util.FromTraceIDContext(ctx), util.FromUserIDContext(ctx))
}

// Operate 操作日志
func Operate(traceID string, userID ...string) *logrus.Entry {
	return logger().Operate(traceID, userID...)
}

// LoginWithContext 登录(登出)日志
func LoginWithContext(ctx context.Context) *logrus.Entry {
	return Login(util.FromTraceIDContext(ctx), util.FromUserIDContext(ctx))
}

// Login 登录(登出)日志
func Login(traceID string, userID string) *logrus.Entry {
	return logger().Login(traceID, userID)
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

func (a *Logger) typeEntry(traceID, fieldType string, userID ...string) *logrus.Entry {
	fields := logrus.Fields{
		FieldKeyTraceID: traceID,
		FieldKeyType:    fieldType,
	}
	if len(userID) > 0 {
		fields[FieldKeyUserID] = userID[0]
	}
	return a.WithFields(fields)
}

// System 系统日志
func (a *Logger) System(traceID string, userID ...string) *logrus.Entry {
	return a.typeEntry(traceID, "system", userID...)
}

// Access 访问日志
func (a *Logger) Access(traceID string, userID ...string) *logrus.Entry {
	return a.typeEntry(traceID, "access", userID...)
}

// Operate 操作日志
func (a *Logger) Operate(traceID string, userID ...string) *logrus.Entry {
	return a.typeEntry(traceID, "operate", userID...)
}

// Login 登录(登出)日志
func (a *Logger) Login(traceID string, userID string) *logrus.Entry {
	return a.typeEntry(traceID, "login", userID)
}
