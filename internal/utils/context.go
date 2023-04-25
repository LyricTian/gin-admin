package utils

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/pkg/encoding/json"
	"gorm.io/gorm"
)

type (
	traceIDCtx    struct{}
	transCtx      struct{}
	rowLockCtx    struct{}
	userIDCtx     struct{}
	userTokenCtx  struct{}
	isRootUserCtx struct{}
	userCacheCtx  struct{}
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

func NewUserToken(ctx context.Context, userToken string) context.Context {
	return context.WithValue(ctx, userTokenCtx{}, userToken)
}

func FromUserToken(ctx context.Context) string {
	v := ctx.Value(userTokenCtx{})
	if v != nil {
		return v.(string)
	}
	return ""
}

func NewIsRootUser(ctx context.Context) context.Context {
	return context.WithValue(ctx, isRootUserCtx{}, true)
}

func FromIsRootUser(ctx context.Context) bool {
	v := ctx.Value(isRootUserCtx{})
	return v != nil && v.(bool)
}

// Set user cache object
type UserCache struct {
	RoleIDs []string `json:"rids"`
}

func ParseUserCache(s string) UserCache {
	var a UserCache
	if s == "" {
		return a
	}

	_ = json.Unmarshal([]byte(s), &a)
	return a
}

func (a UserCache) String() string {
	return json.MarshalToString(a)
}

func NewUserCache(ctx context.Context, userCache UserCache) context.Context {
	return context.WithValue(ctx, userCacheCtx{}, userCache)
}

func FromUserCache(ctx context.Context) UserCache {
	v := ctx.Value(userCacheCtx{})
	if v != nil {
		return v.(UserCache)
	}
	return UserCache{}
}
