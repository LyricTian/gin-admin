package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
)

// GetMenuDB 获取菜单存储
func GetMenuDB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	return getDBWithModel(ctx, defDB, Menu{})
}

// SchemaMenu 菜单对象
type SchemaMenu schema.Menu

// ToMenu 转换为菜单实体
func (a SchemaMenu) ToMenu() *Menu {
	item := &Menu{
		RecordID:   a.RecordID,
		Code:       a.Code,
		Name:       a.Name,
		Type:       a.Type,
		Sequence:   a.Sequence,
		Icon:       a.Icon,
		Path:       a.Path,
		Method:     a.Method,
		ParentID:   a.ParentID,
		ParentPath: a.ParentPath,
		Creator:    a.Creator,
	}
	return item
}

// Menu 菜单实体
type Menu struct {
	Model
	RecordID   string `gorm:"column:record_id;size:36;index;"`    // 记录内码
	Code       string `gorm:"column:code;size:50;"`               // 菜单编号
	Name       string `gorm:"column:name;size:50;index;"`         // 菜单名称
	Type       int    `gorm:"column:type;index;"`                 // 菜单类型(1：模块 2：功能 3：资源)
	Sequence   int    `gorm:"column:sequence;index;"`             // 排序值
	Icon       string `gorm:"column:icon;size:255;"`              // 菜单图标
	Path       string `gorm:"column:path;size:255;"`              // 访问路径
	Method     string `gorm:"column:method;size:50;"`             // 资源请求方式
	ParentID   string `gorm:"column:parent_id;size:36;index;"`    // 父级内码
	ParentPath string `gorm:"column:parent_path;size:518;index;"` // 父级路径
	Creator    string `gorm:"column:creator;size:36;"`            // 创建人
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
		Code:       a.Code,
		Name:       a.Name,
		Type:       a.Type,
		Sequence:   a.Sequence,
		Icon:       a.Icon,
		Path:       a.Path,
		Method:     a.Method,
		ParentID:   a.ParentID,
		ParentPath: a.ParentPath,
		Creator:    a.Creator,
		CreatedAt:  a.CreatedAt,
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
