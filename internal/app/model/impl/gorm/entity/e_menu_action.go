package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetMenuActionDB 菜单动作
func GetMenuActionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, new(MenuAction))
}

// SchemaMenuAction 菜单动作
type SchemaMenuAction schema.MenuAction

// ToMenuAction 转换为菜单动作实体
func (a SchemaMenuAction) ToMenuAction() *MenuAction {
	item := new(MenuAction)
	util.StructMapToStruct(a, item)
	return item
}

// MenuAction 菜单动作实体
type MenuAction struct {
	Model
	MenuID string `gorm:"column:menu_id;size:36;index;default:'';not null;"` // 菜单ID
	Code   string `gorm:"column:code;size:100;default:'';not null;"`         // 动作编号
	Name   string `gorm:"column:name;size:100;default:'';not null;"`         // 动作名称
}

func (a MenuAction) String() string {
	return toString(a)
}

// TableName 表名
func (a MenuAction) TableName() string {
	return a.Model.TableName("menu_action")
}

// ToSchemaMenuAction 转换为菜单动作对象
func (a MenuAction) ToSchemaMenuAction() *schema.MenuAction {
	item := new(schema.MenuAction)
	util.StructMapToStruct(a, item)
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
