package entity

import (
	"context"
	"strings"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
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
	list := make([]*RoleMenu, len(a.Menus))
	for i, item := range a.Menus {
		list[i] = &RoleMenu{
			RoleID:   a.RecordID,
			MenuID:   item.MenuID,
			Action:   strings.Join(item.Actions, ","),
			Resource: strings.Join(item.Resources, ","),
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
		CreatedAt: &a.CreatedAt,
		UpdatedAt: &a.UpdatedAt,
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

// SchemaRoleMenu 角色菜单对象
type SchemaRoleMenu schema.RoleMenu

// ToRoleMenu 转换为角色菜单实体
func (a SchemaRoleMenu) ToRoleMenu(roleID string) *RoleMenu {
	return &RoleMenu{
		RoleID:   roleID,
		MenuID:   a.MenuID,
		Action:   strings.Join(a.Actions, ","),
		Resource: strings.Join(a.Resources, ","),
	}
}

// RoleMenu 角色菜单关联实体
type RoleMenu struct {
	Model
	RoleID   string `gorm:"column:role_id;size:36;index;"` // 角色内码
	MenuID   string `gorm:"column:menu_id;size:36;index;"` // 菜单内码
	Action   string `gorm:"column:action;size:2048;"`      // 动作权限(多个以英文逗号分隔)
	Resource string `gorm:"column:resource;size:2048;"`    // 资源权限(多个以英文逗号分隔)
}

// TableName 表名
func (a RoleMenu) TableName() string {
	return a.Model.TableName("role_menu")
}

// ToSchemaRoleMenu 转换为角色菜单对象
func (a RoleMenu) ToSchemaRoleMenu() *schema.RoleMenu {
	return &schema.RoleMenu{
		MenuID:    a.MenuID,
		Actions:   strings.Split(a.Action, ","),
		Resources: strings.Split(a.Resource, ","),
	}
}

// RoleMenus 角色菜单关联实体列表
type RoleMenus []*RoleMenu

// GetByRoleID 根据角色ID获取角色菜单对象列表
func (a RoleMenus) GetByRoleID(roleID string) []*schema.RoleMenu {
	var list []*schema.RoleMenu
	for _, item := range a {
		if item.RoleID == roleID {
			list = append(list, item.ToSchemaRoleMenu())
		}
	}
	return list
}

// ToSchemaRoleMenus 转换为角色菜单对象列表
func (a RoleMenus) ToSchemaRoleMenus() []*schema.RoleMenu {
	list := make([]*schema.RoleMenu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRoleMenu()
	}
	return list
}

// ToMap 转换为键值映射
func (a RoleMenus) ToMap() map[string]*RoleMenu {
	m := make(map[string]*RoleMenu)
	for _, item := range a {
		m[item.MenuID] = item
	}
	return m
}
