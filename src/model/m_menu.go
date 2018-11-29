package model

import (
	"context"

	"github.com/LyricTian/gin-admin/src/schema"
)

// IMenu 菜单管理
type IMenu interface {
	// 查询分页数据
	QueryPage(ctx context.Context, params schema.MenuQueryParam, pageIndex, pageSize uint) (int64, []*schema.MenuQueryResult, error)
	// 查询选择数据
	QuerySelect(ctx context.Context, params schema.MenuSelectQueryParam) ([]*schema.MenuSelectQueryResult, error)
	// Get 查询指定数据
	Get(ctx context.Context, recordID string) (*schema.Menu, error)
	// Check 检查数据是否存在
	Check(ctx context.Context, recordID string) (bool, error)
	// 检查编号是否存在
	CheckCode(ctx context.Context, code string, parentID string) (bool, error)
	// 根据父级查询分级码
	QueryLevelCodesByParentID(parentID string) ([]string, error)
	// 检查子级是否存在
	CheckChild(ctx context.Context, parentID string) (bool, error)
	// Create 创建数据
	Create(ctx context.Context, item *schema.Menu) error
	// Update 更新数据
	Update(ctx context.Context, recordID string, info map[string]interface{}) error
	// 更新数据
	UpdateWithLevelCode(ctx context.Context, recordID string, info map[string]interface{}, oldLevelCode, newLevelCode string) error
	// Delete 删除数据
	Delete(ctx context.Context, recordID string) error
}
