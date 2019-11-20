package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
)

// IMenuActionResource 菜单动作关联资源管理存储接口
type IMenuActionResource interface {
	// 查询数据
	Query(ctx context.Context, params schema.MenuActionResourceQueryParam, opts ...schema.MenuActionResourceQueryOptions) (*schema.MenuActionResourceQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.MenuActionResourceQueryOptions) (*schema.MenuActionResource, error)
	// 创建数据
	Create(ctx context.Context, item schema.MenuActionResource) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.MenuActionResource) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 根据动作ID删除数据
	DeleteByActionID(ctx context.Context, actionID string) error
	// 根据菜单ID删除数据
	DeleteByMenuID(ctx context.Context, menuID string) error
}
