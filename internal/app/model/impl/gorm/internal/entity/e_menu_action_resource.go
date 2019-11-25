package entity

import (
	"context"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// GetMenuActionResourceDB 菜单动作关联资源
func GetMenuActionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, MenuActionResource{})
}

// SchemaMenuActionResource 菜单动作关联资源
type SchemaMenuActionResource schema.MenuActionResource

// ToMenuActionResource 转换为菜单动作关联资源实体
func (a SchemaMenuActionResource) ToMenuActionResource() *MenuActionResource {
	item := &MenuActionResource{
		RecordID: &a.RecordID,
		ActionID: &a.ActionID,
		Method:   &a.Method,
		Path:     &a.Path,
	}
	return item
}

// MenuActionResource 菜单动作关联资源实体
type MenuActionResource struct {
	Model
	RecordID *string `gorm:"column:record_id;size:36;index;"` // 记录ID
	ActionID *string `gorm:"column:action_id;size:36;index;"` // 菜单动作ID
	Method   *string `gorm:"column:method;size:100;"`         // 资源请求方式(支持正则)
	Path     *string `gorm:"column:path;size:100;"`           // 资源请求路径（支持/:id匹配）
}

func (a MenuActionResource) String() string {
	return toString(a)
}

// TableName 表名
func (a MenuActionResource) TableName() string {
	return a.Model.TableName("menu_action_resource")
}

// ToSchemaMenuActionResource 转换为菜单动作关联资源对象
func (a MenuActionResource) ToSchemaMenuActionResource() *schema.MenuActionResource {
	item := &schema.MenuActionResource{
		RecordID: *a.RecordID,
		ActionID: *a.ActionID,
		Method:   *a.Method,
		Path:     *a.Path,
	}
	return item
}

// MenuActionResources 菜单动作关联资源列表
type MenuActionResources []*MenuActionResource

// ToSchemaMenuActionResources 转换为菜单动作关联资源对象列表
func (a MenuActionResources) ToSchemaMenuActionResources() []*schema.MenuActionResource {
	list := make([]*schema.MenuActionResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuActionResource()
	}
	return list
}
