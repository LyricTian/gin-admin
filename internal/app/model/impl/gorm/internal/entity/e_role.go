package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetRoleDB 获取角色存储
func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, new(Role))
}

// SchemaRole 角色对象
type SchemaRole schema.Role

// ToRole 转换为角色实体
func (a SchemaRole) ToRole() *Role {
	item := new(Role)
	util.StructMapToStruct(a, item)
	return item
}

// Role 角色实体
type Role struct {
	Model
	RecordID string  `gorm:"column:record_id;size:36;index;"` // 记录内码
	Name     *string `gorm:"column:name;size:100;index;"`     // 角色名称
	Sequence *int    `gorm:"column:sequence;index;"`          // 排序值
	Memo     *string `gorm:"column:memo;size:1024;"`          // 备注
	Status   *int    `gorm:"column:status;index;"`            // 状态(1:启用 2:禁用)
	Creator  *string `gorm:"column:creator;size:36;"`         // 创建者
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
	item := new(schema.Role)
	util.StructMapToStruct(a, item)
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

// ----------------------------------------RoleMenu--------------------------------------

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
	RecordID string  `gorm:"column:record_id;size:36;index;"` // 记录ID
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
