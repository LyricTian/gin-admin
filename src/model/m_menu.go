package model

import "gin-admin/src/schema"

// IMenu 菜单管理
type IMenu interface {
	// 查询分页数据
	QueryPage(param schema.MenuQueryParam, pageIndex, pageSize uint) (int64, []*schema.MenuQueryResult, error)
	// 查询选择数据
	QuerySelect(param schema.MenuSelectQueryParam) ([]*schema.MenuSelectQueryResult, error)
	// Get 查询指定数据
	Get(recordID string) (*schema.Menu, error)
	// Create 创建数据
	Create(item *schema.Menu) error
	// Update 更新数据
	Update(recordID string, info map[string]interface{}) error
	// Delete 删除数据
	Delete(recordID string) error
}
