package menu

import (
	"context"

	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/util/structure"
)

// GetMenuActionDB 菜单动作
func GetMenuActionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(MenuAction))
}

// SchemaMenuAction 菜单动作
type SchemaMenuAction schema.MenuAction

// ToMenuAction 转换为菜单动作实体
func (a SchemaMenuAction) ToMenuAction() *MenuAction {
	item := new(MenuAction)
	structure.Copy(a, item)
	return item
}

// MenuAction 菜单动作实体
type MenuAction struct {
	util.Model
	MenuID uint64 `gorm:"index;not null;"` // 菜单ID
	Code   string `gorm:"size:100;"`       // 动作编号
	Name   string `gorm:"size:100;"`       // 动作名称
}

// ToSchemaMenuAction 转换为菜单动作对象
func (a MenuAction) ToSchemaMenuAction() *schema.MenuAction {
	item := new(schema.MenuAction)
	structure.Copy(a, item)
	return item
}

// MenuActions 菜单动作列表
type MenuActions []*MenuAction

// ToSchemaMenuActions 转换为菜单动作对象列表
func (a MenuActions) ToSchemaMenuActions() []*schema.MenuAction {
	list := make([]*schema.MenuAction, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuAction()
	}
	return list
}
