package model

import (
	"context"

	"github.com/LyricTian/gin-admin/src/schema"
)

// IDemo demo存储接口
type IDemo interface {
	// 查询数据
	Query(ctx context.Context, params schema.DemoQueryParam, pp *schema.PaginationParam) ([]*schema.Demo, *schema.PaginationResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string) (*schema.Demo, error)
	// 检查编号是否存在
	CheckCode(ctx context.Context, code string) (bool, error)
	// 创建数据
	Create(ctx context.Context, item schema.Demo) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.Demo) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
