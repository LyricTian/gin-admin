package logger

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/src/util"
	"github.com/sirupsen/logrus"
)

// 定义键名
const (
	StartTimeKey    = "start_time"
	TraceIDKey      = "trace_id"
	UserIDKey       = "user_id"
	SpanIDKey       = "span_id"
	SpanTitleKey    = "span_title"
	SpanFunctionKey = "span_function"
	VersionKey      = "version"
)

// Entry 定义别名
type Entry = logrus.Entry

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
		StartTimeKey:    time.Now(),
		UserIDKey:       FromUserIDContext(ctx),
		TraceIDKey:      FromTraceIDContext(ctx),
		SpanIDKey:       FromSpanIDContext(ctx),
		SpanTitleKey:    title,
		SpanFunctionKey: function,
	}

	return logrus.WithFields(fields)
}
