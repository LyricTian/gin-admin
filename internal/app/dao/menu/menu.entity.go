package menu

import (
	"context"

	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/util/structure"
)

// GetMenuDB 获取菜单存储
func GetMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Menu))
}

// SchemaMenu 菜单对象
type SchemaMenu schema.Menu

// ToMenu 转换为菜单实体
func (a SchemaMenu) ToMenu() *Menu {
	item := new(Menu)
	structure.Copy(a, item)
	return item
}

// Menu 菜单实体
type Menu struct {
	util.Model
	Name       string  `gorm:"size:50;index;default:'';not null;"` // 菜单名称
	Icon       *string `gorm:"size:255;"`                          // 菜单图标
	Router     *string `gorm:"size:255;"`                          // 访问路由
	ParentID   *uint64 `gorm:"index;default:0;"`                   // 父级内码
	ParentPath *string `gorm:"size:512;index;default:'';"`         // 父级路径
	IsShow     int     `gorm:"index;default:0;"`                   // 是否显示(1:显示 2:隐藏)
	Status     int     `gorm:"index;default:0;"`                   // 状态(1:启用 2:禁用)
	Sequence   int     `gorm:"index;default:0;"`                   // 排序值
	Memo       *string `gorm:"size:1024;"`                         // 备注
	Creator    uint64  `gorm:""`                                   // 创建人
}

// ToSchemaMenu 转换为菜单对象
func (a Menu) ToSchemaMenu() *schema.Menu {
	item := new(schema.Menu)
	structure.Copy(a, item)
	return item
}

// Menus 菜单实体列表
type Menus []*Menu

// ToSchemaMenus 转换为菜单对象列表
func (a Menus) ToSchemaMenus() []*schema.Menu {
	list := make([]*schema.Menu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenu()
	}
	return list
}
