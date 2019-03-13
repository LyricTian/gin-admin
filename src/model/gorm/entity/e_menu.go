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

// GetMenuResourceDB 获取菜单资源存储
func GetMenuResourceDB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	return getDBWithModel(ctx, defDB, MenuResource{})
}

// SchemaMenu 菜单对象
type SchemaMenu schema.Menu

// ToMenu 转换为菜单实体
func (a SchemaMenu) ToMenu() *Menu {
	item := &Menu{
		RecordID:   a.RecordID,
		Name:       a.Name,
		Sequence:   a.Sequence,
		Icon:       a.Icon,
		Router:     a.Router,
		Hidden:     a.Hidden,
		ParentID:   a.ParentID,
		ParentPath: a.ParentPath,
		Creator:    a.Creator,
	}
	return item
}

// ToMenuResources 转换为菜单资源列表
func (a SchemaMenu) ToMenuResources() []*MenuResource {
	list := make([]*MenuResource, len(a.Resources))
	for i, item := range a.Resources {
		list[i] = SchemaMenuResource(*item).ToMenuResource(a.RecordID)
	}
	return list
}

// Menu 菜单实体
type Menu struct {
	Model
	RecordID   string `gorm:"column:record_id;size:36;index;"`    // 记录内码
	Name       string `gorm:"column:name;size:50;index;"`         // 菜单名称
	Sequence   int    `gorm:"column:sequence;index;"`             // 排序值
	Icon       string `gorm:"column:icon;size:255;"`              // 菜单图标
	Router     string `gorm:"column:router;size:255;"`            // 访问路由
	Hidden     int    `gorm:"column:hidden;index;"`               // 隐藏菜单(0:不隐藏 1:隐藏)
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
		Name:       a.Name,
		Sequence:   a.Sequence,
		Icon:       a.Icon,
		Router:     a.Router,
		Hidden:     a.Hidden,
		ParentID:   a.ParentID,
		ParentPath: a.ParentPath,
		Creator:    a.Creator,
		CreatedAt:  &a.CreatedAt,
		UpdatedAt:  &a.UpdatedAt,
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

// SchemaMenuResource 菜单资源对象
type SchemaMenuResource schema.MenuResource

// ToMenuResource 转换为菜单资源实体
func (a SchemaMenuResource) ToMenuResource(menuID string) *MenuResource {
	return &MenuResource{
		MenuID:   menuID,
		RecordID: a.RecordID,
		Name:     a.Name,
		Method:   a.Method,
		Path:     a.Path,
	}
}

// MenuResource 菜单资源关联实体
type MenuResource struct {
	Model
	MenuID   string `gorm:"column:menu_id;size:36;index;"`   // 菜单ID
	RecordID string `gorm:"column:record_id;size:36;index;"` // 记录内码
	Name     string `gorm:"column:code;size:50;"`            // 资源名称
	Method   string `gorm:"column:method;size:50;"`          // 请求方式
	Path     string `gorm:"column:path;size:255;"`           // 请求路径
}

// TableName 表名
func (a MenuResource) TableName() string {
	return a.Model.TableName("menu_resource")
}

// ToSchemaMenuResource 转换为菜单资源对象
func (a MenuResource) ToSchemaMenuResource() *schema.MenuResource {
	return &schema.MenuResource{
		RecordID: a.RecordID,
		Name:     a.Name,
		Method:   a.Method,
		Path:     a.Path,
	}
}

// MenuResources 菜单资源关联实体列表
type MenuResources []*MenuResource

// ToSchemaMenuResources 转换为菜单资源列表
func (a MenuResources) ToSchemaMenuResources() []*schema.MenuResource {
	list := make([]*schema.MenuResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuResource()
	}
	return list
}
