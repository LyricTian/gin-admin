package model

import (
	"context"
	"gin-admin/src/schema"
)

// IUser 用户管理
type IUser interface {
	// 查询分页数据
	QueryPage(ctx context.Context, params schema.UserQueryParam, pageIndex, pageSize uint) (int64, []*schema.UserQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, includeRoleIDs bool) (*schema.User, error)
	// 检查用户名
	CheckUserName(ctx context.Context, userName string) (bool, error)
	// 根据用户名查询指定数据
	GetByUserName(ctx context.Context, userName string, includeRoleIDs bool) (*schema.User, error)
	// 创建数据
	Create(ctx context.Context, item *schema.User) error
	// 更新数据
	Update(ctx context.Context, recordID string, info map[string]interface{}) error
	// 更新数据
	UpdateWithRoleIDs(ctx context.Context, recordID string, info map[string]interface{}, roleIDs []string) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
}
