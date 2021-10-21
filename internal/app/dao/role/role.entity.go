package role

import (
	"context"

	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/util/structure"
)

func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Role))
}

type SchemaRole schema.Role

func (a SchemaRole) ToRole() *Role {
	item := new(Role)
	structure.Copy(a, item)
	return item
}

type Role struct {
	util.Model
	Name     string  `gorm:"size:100;index;default:'';not null;"` // 角色名称
	Sequence int     `gorm:"index;default:0;"`                    // 排序值
	Memo     *string `gorm:"size:1024;"`                          // 备注
	Status   int     `gorm:"index;default:0;"`                    // 状态(1:启用 2:禁用)
	Creator  uint64  `gorm:""`                                    // 创建者
}

func (a Role) ToSchemaRole() *schema.Role {
	item := new(schema.Role)
	structure.Copy(a, item)
	return item
}

type Roles []*Role

func (a Roles) ToSchemaRoles() []*schema.Role {
	list := make([]*schema.Role, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}
