package context

import (
	"context"
)

// 定义全局上下文中的键
type (
	transCtx   struct{}
	userIDCtx  struct{}
	traceIDCtx struct{}
)

// NewTrans 创建事务的上下文
func NewTrans(ctx context.Context, trans interface{}) context.Context {
	return context.WithValue(ctx, transCtx{}, trans)
}

// FromTrans 从上下文中获取事务
func FromTrans(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(transCtx{})
	return v, v != nil
}

// NewUserID 创建用户ID的上下文
func NewUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userID)
}

// FromUserID 从上下文中获取用户ID
func FromUserID(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s, s != ""
		}
	}
	return "", false
}

// NewTraceID 创建追踪ID的上下文
func NewTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDCtx{}, traceID)
}

// FromTraceID 从上下文中获取追踪ID
func FromTraceID(ctx context.Context) (string, bool) {
	v := ctx.Value(traceIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s, s != ""
		}
	}
	return "", false
}
