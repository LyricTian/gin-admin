package entity

import (
	"github.com/LyricTian/gin-admin/src/schema"
)

// GetRoleTableName 获取角色表名
func GetRoleTableName() string {
	return Role{}.TableName()
}

// GetRoleMenuTableName 获取角色菜单关联表名
func GetRoleMenuTableName() string {
	return RoleMenu{}.TableName()
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
	}
	return item
}

// Role 角色实体
type Role struct {
	Model
	RecordID string `gorm:"column:record_id;size:36;index;"` // 记录内码
	Name     string `gorm:"column:name;size:100;index;"`            // 角色名称
	Sequence int    `gorm:"column:sequence;index;"`                 // 排序值
	Memo     string `gorm:"column:memo;size:200;"`                  // 备注
	Creator  string `gorm:"column:creator;size:36;"`                // 创建者
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
		RecordID: a.RecordID,
		Name:     a.Name,
		Sequence: a.Sequence,
		Memo:     a.Memo,
	}
	return item
}

// Roles 角色列表
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

// RoleMenus 角色菜单关联列表
type RoleMenus []*RoleMenu

// ToMenuIDs 转换为菜单ID列表
func (a RoleMenus) ToMenuIDs() []string {
	menuIDs := make([]string, len(a))
	for i, item := range a {
		menuIDs[i] = item.MenuID
	}
	return menuIDs
}
