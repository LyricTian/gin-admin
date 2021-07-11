package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/v7/internal/app/schema"
	"github.com/LyricTian/gin-admin/v7/pkg/util/structure"
	"github.com/jinzhu/gorm"
)

// GetMenuActionResourceDB 菜单动作关联资源
func GetMenuActionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(MenuActionResource))
}

// SchemaMenuActionResource 菜单动作关联资源
type SchemaMenuActionResource schema.MenuActionResource

// ToMenuActionResource 转换为菜单动作关联资源实体
func (a SchemaMenuActionResource) ToMenuActionResource() *MenuActionResource {
	item := new(MenuActionResource)
	structure.Copy(a, item)
	return item
}

// MenuActionResource 菜单动作关联资源实体
type MenuActionResource struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"` // 菜单动作ID
	Method   string `gorm:"column:method;size:100;default:'';not null;"`         // 资源请求方式(支持正则)
	Path     string `gorm:"column:path;size:100;default:'';not null;"`           // 资源请求路径（支持/:id匹配）
}

// ToSchemaMenuActionResource 转换为菜单动作关联资源对象
func (a MenuActionResource) ToSchemaMenuActionResource() *schema.MenuActionResource {
	item := new(schema.MenuActionResource)
	structure.Copy(a, item)
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
