package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
)

// IMenuAction 菜单动作管理存储接口
type IMenuAction interface {
	// 查询数据
	Query(ctx context.Context, params schema.MenuActionQueryParam, opts ...schema.MenuActionQueryOptions) (*schema.MenuActionQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.MenuActionQueryOptions) (*schema.MenuAction, error)
	// 创建数据
	Create(ctx context.Context, item schema.MenuAction) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.MenuAction) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 根据菜单ID删除数据
	DeleteByMenuID(ctx context.Context, menuID string) error
}
