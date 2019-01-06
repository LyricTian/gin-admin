package logger

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/LyricTian/gin-admin/src/util"
	"github.com/sirupsen/logrus"
)

// 定义键名
const (
	StartedAtKey    = "started_at"
	TraceIDKey      = "trace_id"
	UserIDKey       = "user_id"
	SpanIDKey       = "span_id"
	SpanTitleKey    = "span_title"
	SpanFunctionKey = "span_function"
	VersionKey      = "version"
)

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
	return util.MustUUID()
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
	return util.MustUUID()
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

// Start 开始写入日志
func Start(ctx context.Context) *Entry {
	return StartSpan(ctx, "", "")
}

// StartSpan 开始一个追踪单元
func StartSpan(ctx context.Context, title, function string) *Entry {
	fields := map[string]interface{}{
		StartedAtKey:    time.Now(),
		UserIDKey:       FromUserIDContext(ctx),
		TraceIDKey:      FromTraceIDContext(ctx),
		SpanIDKey:       FromSpanIDContext(ctx),
		SpanTitleKey:    title,
		SpanFunctionKey: function,
	}

	return newEntry(logrus.WithFields(fields))
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
		e.entry.Info()
	}
}

// WithFields 结构化字段写入
func (e *Entry) WithFields(fields map[string]interface{}) *Entry {
	return newEntry(e.entry.WithFields(fields))
}

// WithField 结构化字段写入
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(map[string]interface{}{key: value})
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

func (e *Entry) done() {
	atomic.StoreInt32(&e.finish, 1)
}
