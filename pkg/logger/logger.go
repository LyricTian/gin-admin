package logger

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Logger Logrus
type Logger = logrus.Logger

// Entry logrus.Entry alias
type Entry = logrus.Entry

// Hook logrus.Hook alias
type Hook = logrus.Hook

type Level = logrus.Level

// Define logger level
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// SetLevel Set logger level
func SetLevel(level Level) {
	logrus.SetLevel(level)
}

// SetFormatter Set logger output format (json/text)
func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}
}

// AddHook Add logger hook
func AddHook(hook Hook) {
	logrus.AddHook(hook)
}

// Define key
const (
	TraceIDKey  = "trace_id"
	UserIDKey   = "user_id"
	UserNameKey = "user_name"
	TagKey      = "tag"
	StackKey    = "stack"
)

type (
	traceIDKey  struct{}
	userIDKey   struct{}
	userNameKey struct{}
	tagKey      struct{}
	stackKey    struct{}
)

func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewUserIDContext(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func FromUserIDContext(ctx context.Context) uint64 {
	v := ctx.Value(userIDKey{})
	if v != nil {
		if s, ok := v.(uint64); ok {
			return s
		}
	}
	return 0
}

func NewUserNameContext(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, userNameKey{}, userName)
}

func FromUserNameContext(ctx context.Context) string {
	v := ctx.Value(userNameKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewTagContext(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, tagKey{}, tag)
}

func FromTagContext(ctx context.Context) string {
	v := ctx.Value(tagKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewStackContext(ctx context.Context, stack error) context.Context {
	return context.WithValue(ctx, stackKey{}, stack)
}

func FromStackContext(ctx context.Context) error {
	v := ctx.Value(stackKey{})
	if v != nil {
		if s, ok := v.(error); ok {
			return s
		}
	}
	return nil
}

// WithContext Use context create entry
func WithContext(ctx context.Context) *Entry {
	fields := logrus.Fields{}

	if v := FromTraceIDContext(ctx); v != "" {
		fields[TraceIDKey] = v
	}

	if v := FromUserIDContext(ctx); v != 0 {
		fields[UserIDKey] = v
	}

	if v := FromUserNameContext(ctx); v != "" {
		fields[UserNameKey] = v
	}

	if v := FromTagContext(ctx); v != "" {
		fields[TagKey] = v
	}

	if v := FromStackContext(ctx); v != nil {
		fields[StackKey] = fmt.Sprintf("%+v", v)
	}

	return logrus.WithContext(ctx).WithFields(fields)
}

// Define logrus alias
var (
	Tracef          = logrus.Tracef
	Debugf          = logrus.Debugf
	Infof           = logrus.Infof
	Warnf           = logrus.Warnf
	Errorf          = logrus.Errorf
	Fatalf          = logrus.Fatalf
	Panicf          = logrus.Panicf
	Printf          = logrus.Printf
	SetOutput       = logrus.SetOutput
	SetReportCaller = logrus.SetReportCaller
	StandardLogger  = logrus.StandardLogger
	ParseLevel      = logrus.ParseLevel
)
