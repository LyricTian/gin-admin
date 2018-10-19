package model

import (
	"context"
	"gin-admin/src/schema"
)

// IMenu 菜单管理
type IMenu interface {
	// 查询分页数据
	QueryPage(ctx context.Context, param schema.MenuQueryParam, pageIndex, pageSize uint) (int64, []*schema.MenuQueryResult, error)
	// 查询选择数据
	QuerySelect(ctx context.Context, param schema.MenuSelectQueryParam) ([]*schema.MenuSelectQueryResult, error)
	// Get 查询指定数据
	Get(ctx context.Context, recordID string) (*schema.Menu, error)
	// Create 创建数据
	Create(ctx context.Context, item *schema.Menu) error
	// Update 更新数据
	Update(ctx context.Context, recordID string, info map[string]interface{}) error
	// Delete 删除数据
	Delete(ctx context.Context, recordID string) error
}
