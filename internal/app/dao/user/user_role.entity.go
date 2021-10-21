package user

import (
	"context"

	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/util/structure"
)

func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(UserRole))
}

type SchemaUserRole schema.UserRole

func (a SchemaUserRole) ToUserRole() *UserRole {
	item := new(UserRole)
	structure.Copy(a, item)
	return item
}

type UserRole struct {
	util.Model
	UserID uint64 `gorm:"index;default:0;"` // 用户内码
	RoleID uint64 `gorm:"index;default:0;"` // 角色内码
}

func (a UserRole) ToSchemaUserRole() *schema.UserRole {
	item := new(schema.UserRole)
	structure.Copy(a, item)
	return item
}

type UserRoles []*UserRole

func (a UserRoles) ToSchemaUserRoles() []*schema.UserRole {
	list := make([]*schema.UserRole, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserRole()
	}
	return list
}
