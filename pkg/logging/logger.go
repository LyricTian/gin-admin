package logging

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

const (
	TagKeyMain     = "main"
	TagKeyRecovery = "recovery"
	TagKeyRequest  = "request"
	TagKeyLogin    = "login"
	TagKeyLogout   = "logout"
	TagKeySystem   = "system"
	TagKeyOperate  = "operate"
)

type (
	ctxLoggerKey  struct{}
	ctxTraceIDKey struct{}
	ctxUserIDKey  struct{}
	ctxTagKey     struct{}
	ctxStackKey   struct{}
)

func NewLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

func FromLogger(ctx context.Context) *zap.Logger {
	v := ctx.Value(ctxLoggerKey{})
	if v != nil {
		if vv, ok := v.(*zap.Logger); ok {
			return vv
		}
	}
	return zap.L()
}

func NewTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ctxTraceIDKey{}, traceID)
}

func FromTraceID(ctx context.Context) string {
	v := ctx.Value(ctxTraceIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxUserIDKey{}, userID)
}

func FromUserID(ctx context.Context) string {
	v := ctx.Value(ctxUserIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewTag(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, ctxTagKey{}, tag)
}

func FromTag(ctx context.Context) string {
	v := ctx.Value(ctxTagKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewStack(ctx context.Context, stack string) context.Context {
	return context.WithValue(ctx, ctxStackKey{}, stack)
}

func FromStack(ctx context.Context) string {
	v := ctx.Value(ctxStackKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func Context(ctx context.Context) *zap.Logger {
	var fields []zap.Field
	if v := FromTraceID(ctx); v != "" {
		fields = append(fields, zap.String("trace_id", v))
	}
	if v := FromUserID(ctx); v != "" {
		fields = append(fields, zap.String("user_id", v))
	}
	if v := FromTag(ctx); v != "" {
		fields = append(fields, zap.String("tag", v))
	}
	if v := FromStack(ctx); v != "" {
		fields = append(fields, zap.String("stack", v))
	}
	return FromLogger(ctx).With(fields...)
}

type PrintLogger struct{}

func (a *PrintLogger) Printf(format string, args ...interface{}) {
	zap.L().Info(fmt.Sprintf(format, args...))
}
