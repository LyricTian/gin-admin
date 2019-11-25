package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// GetRoleDB 获取角色存储
func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, Role{})
}

// SchemaRole 角色对象
type SchemaRole schema.Role

// ToRole 转换为角色实体
func (a SchemaRole) ToRole() *Role {
	item := &Role{
		RecordID: a.RecordID,
		Name:     &a.Name,
		Sequence: &a.Sequence,
		Memo:     &a.Memo,
		Creator:  &a.Creator,
	}
	return item
}

// Role 角色实体
type Role struct {
	Model
	RecordID string  `gorm:"column:record_id;size:36;index;"` // 记录内码
	Name     *string `gorm:"column:name;size:100;index;"`     // 角色名称
	Sequence *int    `gorm:"column:sequence;index;"`          // 排序值
	Memo     *string `gorm:"column:memo;size:1024;"`          // 备注
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
	item := &schema.Role{
		RecordID:  a.RecordID,
		Name:      *a.Name,
		Sequence:  *a.Sequence,
		Memo:      *a.Memo,
		Creator:   *a.Creator,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
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
