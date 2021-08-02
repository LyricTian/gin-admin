package user

import (
	"context"

	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/util/structure"
)

// GetUserRoleDB 获取用户角色关联存储
func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(UserRole))
}

// SchemaUserRole 用户角色
type SchemaUserRole schema.UserRole

// ToUserRole 转换为角色菜单实体
func (a SchemaUserRole) ToUserRole() *UserRole {
	item := new(UserRole)
	structure.Copy(a, item)
	return item
}

// UserRole 用户角色关联实体
type UserRole struct {
	util.Model
	UserID uint64 `gorm:"index;default:0;"` // 用户内码
	RoleID uint64 `gorm:"index;default:0;"` // 角色内码
}

// ToSchemaUserRole 转换为用户角色对象
func (a UserRole) ToSchemaUserRole() *schema.UserRole {
	item := new(schema.UserRole)
	structure.Copy(a, item)
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
