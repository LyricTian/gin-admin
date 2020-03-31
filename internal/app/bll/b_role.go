package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
)

// IRole 角色管理业务逻辑接口
type IRole interface {
	// 查询数据
	Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.RoleQueryOptions) (*schema.Role, error)
	// 创建数据
	Create(ctx context.Context, item schema.Role) (*schema.RecordIDResult, error)
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.Role) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
