package model

import (
	"context"

	"github.com/LyricTian/gin-admin/src/schema"
)

// IRole 角色管理
type IRole interface {
	// 查询分页数据
	QueryPage(ctx context.Context, params schema.RolePageQueryParam, pageIndex, pageSize uint) (int, []*schema.Role, error)
	// 查询列表数据
	QueryList(ctx context.Context, params schema.RoleListQueryParam) ([]*schema.Role, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, includeMenuIDs bool) (*schema.Role, error)
	// 检查名称
	CheckName(ctx context.Context, name string) (bool, error)
	// 创建数据
	Create(ctx context.Context, item schema.Role) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.Role) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
