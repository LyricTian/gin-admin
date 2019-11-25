package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
)

// IRoleMenu 角色菜单存储接口
type IRoleMenu interface {
	// 查询数据
	Query(ctx context.Context, params schema.RoleMenuQueryParam, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenuQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error)
	// 创建数据
	Create(ctx context.Context, item schema.RoleMenu) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.RoleMenu) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 根据角色ID删除数据
	DeleteByRoleID(ctx context.Context, roleID string) error
}
