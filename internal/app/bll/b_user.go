package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
)

// IUser 用户管理业务逻辑接口
type IUser interface {
	// 查询数据
	Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error)
	// 查询显示项数据
	QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error)
	// 创建数据
	Create(ctx context.Context, item schema.User) (*schema.RecordIDResult, error)
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.User) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
