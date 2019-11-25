package entity

import (
	"context"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// GetRoleMenuDB 角色菜单
func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, RoleMenu{})
}

// SchemaRoleMenu 角色菜单
type SchemaRoleMenu schema.RoleMenu

// ToRoleMenu 转换为角色菜单实体
func (a SchemaRoleMenu) ToRoleMenu() *RoleMenu {
	item := &RoleMenu{
		RecordID: &a.RecordID,
		RoleID:   &a.RoleID,
		MenuID:   &a.MenuID,
		ActionID: &a.ActionID,
	}
	return item
}

// RoleMenu 角色菜单实体
type RoleMenu struct {
	Model
	RecordID *string `gorm:"column:record_id;size:36;index;"` // 记录ID
	RoleID   *string `gorm:"column:role_id;size:36;index;"`   // 角色ID
	MenuID   *string `gorm:"column:menu_id;size:36;index;"`   // 菜单ID
	ActionID *string `gorm:"column:action_id;size:36;index;"` // 动作ID
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
	item := &schema.RoleMenu{
		RecordID: *a.RecordID,
		RoleID:   *a.RoleID,
		MenuID:   *a.MenuID,
		ActionID: *a.ActionID,
	}
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
