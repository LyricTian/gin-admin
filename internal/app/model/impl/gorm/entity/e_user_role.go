package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetUserRoleDB 获取用户角色关联存储
func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, new(UserRole))
}

// SchemaUserRole 用户角色
type SchemaUserRole schema.UserRole

// ToUserRole 转换为角色菜单实体
func (a SchemaUserRole) ToUserRole() *UserRole {
	item := new(UserRole)
	util.StructMapToStruct(a, item)
	return item
}

// UserRole 用户角色关联实体
type UserRole struct {
	Model
	UserID string `gorm:"column:user_id;size:36;index;default:'';not null;"` // 用户内码
	RoleID string `gorm:"column:role_id;size:36;index;default:'';not null;"` // 角色内码
}

// TableName 表名
func (a UserRole) TableName() string {
	return a.Model.TableName("user_role")
}

// ToSchemaUserRole 转换为用户角色对象
func (a UserRole) ToSchemaUserRole() *schema.UserRole {
	item := new(schema.UserRole)
	util.StructMapToStruct(a, item)
	return item
}

// UserRoles 用户角色关联列表
type UserRoles []*UserRole

// ToSchemaUserRoles 转换为用户角色对象列表
func (a UserRoles) ToSchemaUserRoles() []*schema.UserRole {
	list := make([]*schema.UserRole, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserRole()
	}
	return list
}
