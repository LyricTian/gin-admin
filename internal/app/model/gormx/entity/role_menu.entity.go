package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/v7/internal/app/schema"
	"github.com/LyricTian/gin-admin/v7/pkg/util/structure"
	"github.com/jinzhu/gorm"
)

// GetRoleMenuDB 角色菜单
func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(RoleMenu))
}

// SchemaRoleMenu 角色菜单
type SchemaRoleMenu schema.RoleMenu

// ToRoleMenu 转换为角色菜单实体
func (a SchemaRoleMenu) ToRoleMenu() *RoleMenu {
	item := new(RoleMenu)
	structure.Copy(a, item)
	return item
}

// RoleMenu 角色菜单实体
type RoleMenu struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	RoleID   string `gorm:"column:role_id;size:36;index;default:'';not null;"`   // 角色ID
	MenuID   string `gorm:"column:menu_id;size:36;index;default:'';not null;"`   // 菜单ID
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"` // 动作ID
}

// ToSchemaRoleMenu 转换为角色菜单对象
func (a RoleMenu) ToSchemaRoleMenu() *schema.RoleMenu {
	item := new(schema.RoleMenu)
	structure.Copy(a, item)
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
