package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/v6/internal/app/schema"
)

// IMenu 菜单管理业务逻辑接口
type IMenu interface {
	// 初始化菜单数据
	InitData(ctx context.Context, dataFile string) error
	// 查询数据
	Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error)
	// 创建数据
	Create(ctx context.Context, item schema.Menu) (*schema.IDResult, error)
	// 更新数据
	Update(ctx context.Context, id string, item schema.Menu) error
	// 删除数据
	Delete(ctx context.Context, id string) error
	// 更新状态
	UpdateStatus(ctx context.Context, id string, status int) error
}
