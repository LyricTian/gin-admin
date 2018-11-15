package model

import (
	"context"
	"gin-admin/src/schema"
)

// IDemo 示例程序
type IDemo interface {
	// 查询分页数据
	QueryPage(ctx context.Context, param schema.DemoQueryParam, pageIndex, pageSize uint) (int64, []*schema.DemoQueryResult, error)
	// Get 查询指定数据
	Get(ctx context.Context, recordID string) (*schema.Demo, error)
	// Check 检查数据是否存在
	Check(ctx context.Context, recordID string) (bool, error)
	// Create 创建数据
	Create(ctx context.Context, item *schema.Demo) error
	// Update 更新数据
	Update(ctx context.Context, recordID string, info map[string]interface{}) error
	// Delete 删除数据
	Delete(ctx context.Context, recordID string) error
}
