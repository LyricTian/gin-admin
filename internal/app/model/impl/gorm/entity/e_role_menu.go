package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetRoleMenuDB 角色菜单
func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, new(RoleMenu))
}

// SchemaRoleMenu 角色菜单
type SchemaRoleMenu schema.RoleMenu

// ToRoleMenu 转换为角色菜单实体
func (a SchemaRoleMenu) ToRoleMenu() *RoleMenu {
	item := new(RoleMenu)
	util.StructMapToStruct(a, item)
	return item
}

// RoleMenu 角色菜单实体
type RoleMenu struct {
	Model
	RoleID   string `gorm:"column:role_id;size:36;index;default:'';not null;"`   // 角色ID
	MenuID   string `gorm:"column:menu_id;size:36;index;default:'';not null;"`   // 菜单ID
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"` // 动作ID
}

func (a RoleMenu) String() string {
	return toString(a)
}

// TableName 表名
func (a RoleMenu) TableName() string {
	return a.Model.TableName("role_menu")
}

// ToSchemaRoleMenu 转换为角色菜单对象
func (a RoleMenu) ToSchemaRoleMenu() *schema.RoleMenu {
	item := new(schema.RoleMenu)
	util.StructMapToStruct(a, item)
	return item
}

// RoleMenus 角色菜单列表
type RoleMenus []*RoleMenu

// ToSchemaRoleMenus 转换为角色菜单对象列表
func (a RoleMenus) ToSchemaRoleMenus() []*schema.RoleMenu {
	list := make([]*schema.RoleMenu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRoleMenu()
	}
	return list
}
