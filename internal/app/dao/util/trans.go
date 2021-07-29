package util

import (
	"context"

	"github.com/LyricTian/gin-admin/v8/internal/app/contextx"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// TransSet 注入Trans
var TransSet = wire.NewSet(wire.Struct(new(Trans), "*"))

// Trans 事务管理
type Trans struct {
	DB *gorm.DB
}

// Exec 执行事务
func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := contextx.FromTrans(ctx); ok {
		return fn(ctx)
	}

	return a.DB.Transaction(func(db *gorm.DB) error {
		return fn(contextx.NewTrans(ctx, db))
	})
}
