package menu

import (
	"context"

	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/util/structure"
)

// GetMenuActionResourceDB 菜单动作关联资源
func GetMenuActionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(MenuActionResource))
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
	util.Model
	ActionID uint64 `gorm:"index;not null;"` // 菜单动作ID
	Method   string `gorm:"size:50;"`        // 资源请求方式(支持正则)
	Path     string `gorm:"size:255;"`       // 资源请求路径（支持/:id匹配）
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
