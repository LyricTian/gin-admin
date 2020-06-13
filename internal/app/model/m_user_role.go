package model

import (
	"context"

	"github.com/LyricTian/gin-admin/v6/internal/app/schema"
)

// IUserRole 用户角色存储接口
type IUserRole interface {
	// 查询数据
	Query(ctx context.Context, params schema.UserRoleQueryParam, opts ...schema.UserRoleQueryOptions) (*schema.UserRoleQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, id string, opts ...schema.UserRoleQueryOptions) (*schema.UserRole, error)
	// 创建数据
	Create(ctx context.Context, item schema.UserRole) error
	// 更新数据
	Update(ctx context.Context, id string, item schema.UserRole) error
	// 删除数据
	Delete(ctx context.Context, id string) error
	// 根据用户ID删除数据
	DeleteByUserID(ctx context.Context, userID string) error
}
