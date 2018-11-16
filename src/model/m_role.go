package model

import (
	"context"
	"gin-admin/src/schema"
)

// IRole 角色管理
type IRole interface {
	// 查询分页数据
	QueryPage(ctx context.Context, params schema.RoleQueryParam, pageIndex, pageSize uint) (int64, []*schema.RoleQueryResult, error)
	// 查询选择数据
	QuerySelect(ctx context.Context, params schema.RoleSelectQueryParam) ([]*schema.RoleSelectQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, includeMenuIDs bool) (*schema.Role, error)
	// Check 检查数据是否存在
	Check(ctx context.Context, recordID string) (bool, error)
	// 检查名称
	CheckName(ctx context.Context, name string) (bool, error)
	// 创建数据
	Create(ctx context.Context, item *schema.Role) error
	// 更新数据
	Update(ctx context.Context, recordID string, info map[string]interface{}) error
	// 更新数据
	UpdateWithMenuIDs(ctx context.Context, recordID string, info map[string]interface{}, menuIDs []string) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
}
