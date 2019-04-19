package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
)

// IUser 用户对象存储接口
type IUser interface {
	// 查询数据
	Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error)
	// 创建数据
	Create(ctx context.Context, item schema.User) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.User) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
	// 更新密码
	UpdatePassword(ctx context.Context, recordID, password string) error
}
