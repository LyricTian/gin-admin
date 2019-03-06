package logger

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// 定义键名
const (
	StartedAtKey     = "started_at"
	TraceIDKey       = "trace_id"
	UserIDKey        = "user_id"
	SpanIDKey        = "span_id"
	SpanTitleKey     = "span_title"
	SpanFunctionKey  = "span_function"
	VersionKey       = "version"
	TimeConsumingKey = "time_consuming"
)

// TraceIDFunc 定义获取跟踪ID的函数
type TraceIDFunc func() string

var (
	version     string
	traceIDFunc TraceIDFunc
)

// Logger 定义日志别名
type Logger = logrus.Logger

// Hook 定义日志钩子别名
type Hook = logrus.Hook

// StandardLogger 获取标准日志
func StandardLogger() *Logger {
	return logrus.StandardLogger()
}

// SetLevel 设定日志级别
func SetLevel(level int) {
	logrus.SetLevel(logrus.Level(level))
}

// SetFormatter 设定日志输出格式
func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
}

// SetOutput 设定日志输出
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// SetVersion 设定版本
func SetVersion(v string) {
	version = v
}

// SetTraceIDFunc 设定追踪ID的处理函数
func SetTraceIDFunc(fn TraceIDFunc) {
	traceIDFunc = fn
}

// AddHook 增加日志钩子
func AddHook(hook Hook) {
	logrus.AddHook(hook)
}

func getTraceID() string {
	if traceIDFunc != nil {
		return traceIDFunc()
	}
	return time.Now().Format("2006.01.02.15.04.05.000")
}

type (
	traceIDContextKey struct{}
	spanIDContextKey  struct{}
	userIDContextKey  struct{}
)

// NewTraceIDContext 创建跟踪ID上下文
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDContextKey{}, traceID)
}

// FromTraceIDContext 从上下文中获取跟踪ID
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return getTraceID()
}

// NewSpanIDContext 创建跟踪单元ID上下文
func NewSpanIDContext(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanIDContextKey{}, spanID)
}

// FromSpanIDContext 从上下文中获取跟踪单元ID
func FromSpanIDContext(ctx context.Context) string {
	v := ctx.Value(spanIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return getTraceID()
}

// NewUserIDContext 创建用户ID上下文
func NewUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

// FromUserIDContext 从上下文中获取用户ID
func FromUserIDContext(ctx context.Context) string {
	v := ctx.Value(userIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// StartSpan 开始一个追踪单元
func StartSpan(ctx context.Context, title, funcName string) *Entry {
	fields := map[string]interface{}{
		StartedAtKey:    time.Now(),
		UserIDKey:       FromUserIDContext(ctx),
		TraceIDKey:      FromTraceIDContext(ctx),
		SpanIDKey:       FromSpanIDContext(ctx),
		SpanTitleKey:    title,
		SpanFunctionKey: funcName,
		VersionKey:      version,
	}

	return newEntry(logrus.WithFields(fields))
}

// StartSpanWithCall 开始一个追踪单元（回调执行）
func StartSpanWithCall(ctx context.Context, title, funcName string) func() *Entry {
	return func() *Entry {
		return StartSpan(ctx, title, funcName)
	}
}

func newEntry(entry *logrus.Entry) *Entry {
	return &Entry{entry: entry}
}

// Entry 定义统一的日志写入方式
type Entry struct {
	entry  *logrus.Entry
	finish int32
}

// Finish 完成，如果没有触发写入则手动触发Info级别的日志写入
func (e *Entry) Finish() {
	if atomic.CompareAndSwapInt32(&e.finish, 0, 1) {
		e.done()
		e.entry.Info()
	}
}

func (e *Entry) checkAndDelete(fields map[string]interface{}, keys ...string) {
	for _, key := range keys {
		if _, ok := fields[key]; ok {
			delete(fields, key)
		}
	}
}

// WithFields 结构化字段写入
func (e *Entry) WithFields(fields map[string]interface{}) *Entry {
	e.checkAndDelete(fields,
		StartedAtKey,
		TraceIDKey,
		UserIDKey,
		SpanIDKey,
		SpanTitleKey,
		SpanFunctionKey,
		VersionKey,
		TimeConsumingKey)
	return newEntry(e.entry.WithFields(fields))
}

// WithField 结构化字段写入
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(map[string]interface{}{key: value})
}

// Fatalf 重大错误日志
func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.done()
	e.entry.Fatalf(format, args...)
}

// Errorf 错误日志
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.done()
	e.entry.Errorf(format, args...)
}

// Warnf 警告日志
func (e *Entry) Warnf(format string, args ...interface{}) {
	e.done()
	e.entry.Warnf(format, args...)
}

// Infof 消息日志
func (e *Entry) Infof(format string, args ...interface{}) {
	e.done()
	e.entry.Infof(format, args...)
}

// Printf 消息日志
func (e *Entry) Printf(format string, args ...interface{}) {
	e.done()
	e.entry.Printf(format, args...)
}

// Debugf 写入调试日志
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.done()
	e.entry.Debugf(format, args...)
}

func (e *Entry) copyEntry(entry *logrus.Entry) *logrus.Entry {
	newEntry := logrus.NewEntry(entry.Logger)
	newEntry.Data = make(logrus.Fields)
	newEntry.Time = entry.Time
	newEntry.Level = entry.Level
	newEntry.Message = entry.Message
	for k, v := range entry.Data {
		newEntry.Data[k] = v
	}
	return newEntry
}

func (e *Entry) done() {
	entry := e.copyEntry(e.entry)
	entry.Time = time.Now()
	if v, ok := entry.Data[StartedAtKey]; ok {
		if startedAt, ok := v.(time.Time); ok {
			entry.Data[TimeConsumingKey] = entry.Time.Sub(startedAt).Nanoseconds() / 1e3
			delete(entry.Data, StartedAtKey)
		}
	}
	e.entry = entry
	atomic.CompareAndSwapInt32(&e.finish, 0, 1)
}
