package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
)

// GetRoleDB 获取角色存储
func GetRoleDB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	return getDBWithModel(ctx, defDB, Role{})
}

// GetRoleMenuDB 获取角色菜单关联存储
func GetRoleMenuDB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	return getDBWithModel(ctx, defDB, RoleMenu{})
}

// SchemaRole 角色对象
type SchemaRole schema.Role

// ToRole 转换为角色实体
func (a SchemaRole) ToRole() *Role {
	item := &Role{
		RecordID: a.RecordID,
		Name:     a.Name,
		Sequence: a.Sequence,
		Memo:     a.Memo,
		Creator:  a.Creator,
	}
	return item
}

// ToRoleMenus 转换为角色菜单实体列表
func (a SchemaRole) ToRoleMenus() []*RoleMenu {
	list := make([]*RoleMenu, len(a.MenuIDs))
	for i, menuID := range a.MenuIDs {
		list[i] = &RoleMenu{
			RoleID: a.RecordID,
			MenuID: menuID,
		}
	}
	return list
}

// Role 角色实体
type Role struct {
	Model
	RecordID string `gorm:"column:record_id;size:36;index;"` // 记录内码
	Name     string `gorm:"column:name;size:100;index;"`     // 角色名称
	Sequence int    `gorm:"column:sequence;index;"`          // 排序值
	Memo     string `gorm:"column:memo;size:200;"`           // 备注
	Creator  string `gorm:"column:creator;size:36;"`         // 创建者
}

func (a Role) String() string {
	return toString(a)
}

// TableName 表名
func (a Role) TableName() string {
	return a.Model.TableName("role")
}

// ToSchemaRole 转换为角色对象
func (a Role) ToSchemaRole() *schema.Role {
	item := &schema.Role{
		RecordID:  a.RecordID,
		Name:      a.Name,
		Sequence:  a.Sequence,
		Memo:      a.Memo,
		Creator:   a.Creator,
		CreatedAt: a.CreatedAt,
	}
	return item
}

// Roles 角色实体列表
type Roles []*Role

// ToSchemaRoles 转换为角色对象列表
func (a Roles) ToSchemaRoles() []*schema.Role {
	list := make([]*schema.Role, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}

// RoleMenu 角色菜单关联实体
type RoleMenu struct {
	Model
	RoleID string `gorm:"column:role_id;size:36;index;"` // 角色内码
	MenuID string `gorm:"column:menu_id;size:36;index;"` // 菜单内码
}

// TableName 表名
func (a RoleMenu) TableName() string {
	return a.Model.TableName("role_menu")
}

// RoleMenus 角色菜单关联实体列表
type RoleMenus []*RoleMenu

// ToMenuIDs 转换为菜单ID列表
func (a RoleMenus) ToMenuIDs() []string {
	menuIDs := make([]string, len(a))
	for i, item := range a {
		menuIDs[i] = item.MenuID
	}
	return menuIDs
}
