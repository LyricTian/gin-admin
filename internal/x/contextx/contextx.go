package contextx

import (
	"context"

	"gorm.io/gorm"
)

type (
	traceIDCtx struct{}
	transCtx   struct{}
	rowLockCtx struct{}
	userIDCtx  struct{}
	roleIDsCtx struct{}
)

func NewTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDCtx{}, traceID)
}

func FromTraceID(ctx context.Context) string {
	v := ctx.Value(traceIDCtx{})
	if v != nil {
		return v.(string)
	}
	return ""
}

func NewTrans(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, transCtx{}, db)
}

func FromTrans(ctx context.Context) (*gorm.DB, bool) {
	v := ctx.Value(transCtx{})
	if v != nil {
		return v.(*gorm.DB), true
	}
	return nil, false
}

func NewRowLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, rowLockCtx{}, true)
}

func FromRowLock(ctx context.Context) bool {
	v := ctx.Value(rowLockCtx{})
	return v != nil && v.(bool)
}

func NewUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userID)
}

func FromUserID(ctx context.Context) string {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		return v.(string)
	}
	return ""
}

func NewRoleIDs(ctx context.Context, roleIDs []string) context.Context {
	return context.WithValue(ctx, roleIDsCtx{}, roleIDs)
}

func FromRoleIDs(ctx context.Context) []string {
	v := ctx.Value(roleIDsCtx{})
	if v != nil {
		return v.([]string)
	}
	return nil
}
