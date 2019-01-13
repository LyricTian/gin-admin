package model

import (
	"context"

	"github.com/LyricTian/gin-admin/src/schema"
)

// IDemo 示例程序
type IDemo interface {
	// 查询分页数据
	QueryPage(ctx context.Context, params schema.DemoQueryParam, pageIndex, pageSize uint) (int, []schema.DemoQueryResult, error)

	// 查询指定数据
	Get(ctx context.Context, recordID string) (*schema.Demo, error)

	// 检查数据是否存在
	Check(ctx context.Context, recordID string) (bool, error)

	// 检查编号是否存在
	CheckCode(ctx context.Context, code string) (bool, error)

	// 创建数据
	Create(ctx context.Context, item schema.Demo) error

	// 更新数据
	Update(ctx context.Context, recordID string, item schema.Demo) error

	// 删除数据
	Delete(ctx context.Context, trans interface{}, recordID string) error

	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
