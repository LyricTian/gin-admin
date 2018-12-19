package logger

import (
	"context"
	"sync"

	"github.com/LyricTian/gin-admin/src/util"
	"github.com/sirupsen/logrus"
)

// 定义日志中使用的键名
const (
	FieldKeyType    = "type"
	FieldKeyTraceID = "trace_id"
	FieldKeyUserID  = "user_id"
)

var (
	internalLogger *Logger
	once           sync.Once
	defaultOptions = options{
		level:  5,
		format: "text",
	}
)

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

// System 系统日志
func System(ctx context.Context) *logrus.Entry {
	return logger().System(ctx)
}

// Access 访问日志
func Access(ctx context.Context) *logrus.Entry {
	return logger().Access(ctx)
}

// Operate 操作日志
func Operate(ctx context.Context) *logrus.Entry {
	return logger().Operate(ctx)
}

// Login 登录(登出)日志
func Login(ctx context.Context) *logrus.Entry {
	return logger().Login(ctx)
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

func (a *Logger) typeEntry(ctx context.Context, fieldType string) *logrus.Entry {
	fields := logrus.Fields{
		FieldKeyType: fieldType,
	}

	if traceID := util.FromTraceIDContext(ctx); traceID != "" {
		fields[FieldKeyTraceID] = traceID
	}
	if userID := util.FromUserIDContext(ctx); userID != "" {
		fields[FieldKeyUserID] = userID
	}

	return a.WithFields(fields)
}

// System 系统日志
func (a *Logger) System(ctx context.Context) *logrus.Entry {
	return a.typeEntry(ctx, "system")
}

// Access 访问日志
func (a *Logger) Access(ctx context.Context) *logrus.Entry {
	return a.typeEntry(ctx, "access")
}

// Operate 操作日志
func (a *Logger) Operate(ctx context.Context) *logrus.Entry {
	return a.typeEntry(ctx, "operate")
}

// Login 登录(登出)日志
func (a *Logger) Login(ctx context.Context) *logrus.Entry {
	return a.typeEntry(ctx, "login")
}
