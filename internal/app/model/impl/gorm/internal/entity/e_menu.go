package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// GetMenuDB 获取菜单存储
func GetMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, Menu{})
}

// SchemaMenu 菜单对象
type SchemaMenu schema.Menu

// ToMenu 转换为菜单实体
func (a SchemaMenu) ToMenu() *Menu {
	item := &Menu{
		RecordID:   a.RecordID,
		Name:       &a.Name,
		Sequence:   &a.Sequence,
		Icon:       &a.Icon,
		Router:     &a.Router,
		ParentID:   &a.ParentID,
		ParentPath: &a.ParentPath,
		Status:     &a.Status,
		Memo:       &a.Memo,
		Creator:    &a.Creator,
	}
	return item
}

// Menu 菜单实体
type Menu struct {
	Model
	RecordID   string  `gorm:"column:record_id;size:36;index;"`    // 记录内码
	Name       *string `gorm:"column:name;size:50;index;"`         // 菜单名称
	Sequence   *int    `gorm:"column:sequence;index;"`             // 排序值
	Icon       *string `gorm:"column:icon;size:255;"`              // 菜单图标
	Router     *string `gorm:"column:router;size:255;"`            // 访问路由
	ParentID   *string `gorm:"column:parent_id;size:36;index;"`    // 父级内码
	ParentPath *string `gorm:"column:parent_path;size:518;index;"` // 父级路径
	Status     *int    `gorm:"column:status;index;"`               // 状态(1:正常 2:隐藏)
	Memo       *string `gorm:"column:status;size:1024;"`           // 备注
	Creator    *string `gorm:"column:creator;size:36;"`            // 创建人
}

func (a Menu) String() string {
	return toString(a)
}

// TableName 表名
func (a Menu) TableName() string {
	return a.Model.TableName("menu")
}

// ToSchemaMenu 转换为菜单对象
func (a Menu) ToSchemaMenu() *schema.Menu {
	item := &schema.Menu{
		RecordID:   a.RecordID,
		Name:       *a.Name,
		Sequence:   *a.Sequence,
		Icon:       *a.Icon,
		Router:     *a.Router,
		ParentID:   *a.ParentID,
		ParentPath: *a.ParentPath,
		Creator:    *a.Creator,
		Status:     *a.Status,
		Memo:       *a.Memo,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
	return item
}

// Menus 菜单实体列表
type Menus []*Menu

// ToSchemaMenus 转换为菜单对象列表
func (a Menus) ToSchemaMenus() []*schema.Menu {
	list := make([]*schema.Menu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenu()
	}
	return list
}
