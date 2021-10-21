package role

import (
	"context"

	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/util/structure"
)

func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(RoleMenu))
}

type SchemaRoleMenu schema.RoleMenu

func (a SchemaRoleMenu) ToRoleMenu() *RoleMenu {
	item := new(RoleMenu)
	structure.Copy(a, item)
	return item
}

type RoleMenu struct {
	util.Model
	RoleID   uint64 `gorm:"index;not null;"` // 角色ID
	MenuID   uint64 `gorm:"index;not null;"` // 菜单ID
	ActionID uint64 `gorm:"index;not null;"` // 动作ID
}

func (a RoleMenu) ToSchemaRoleMenu() *schema.RoleMenu {
	item := new(schema.RoleMenu)
	structure.Copy(a, item)
	return item
}

type RoleMenus []*RoleMenu

func (a RoleMenus) ToSchemaRoleMenus() []*schema.RoleMenu {
	list := make([]*schema.RoleMenu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRoleMenu()
	}
	return list
}
