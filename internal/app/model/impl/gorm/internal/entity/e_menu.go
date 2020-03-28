package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetMenuDB 获取菜单存储
func GetMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, new(Menu))
}

// SchemaMenu 菜单对象
type SchemaMenu schema.Menu

// ToMenu 转换为菜单实体
func (a SchemaMenu) ToMenu() *Menu {
	item := new(Menu)
	util.StructMapToStruct(a, item)
	return item
}

// Menu 菜单实体
type Menu struct {
	Model
	RecordID   string  `gorm:"column:record_id;size:36;index;not null;"` // 记录内码
	Name       *string `gorm:"column:name;size:50;index;"`               // 菜单名称
	Sequence   *int    `gorm:"column:sequence;index;"`                   // 排序值
	Icon       *string `gorm:"column:icon;size:255;"`                    // 菜单图标
	Router     *string `gorm:"column:router;size:255;"`                  // 访问路由
	ParentID   *string `gorm:"column:parent_id;size:36;index;"`          // 父级内码
	ParentPath *string `gorm:"column:parent_path;size:518;index;"`       // 父级路径
	ShowStatus *int    `gorm:"column:show_status;index;"`                // 状态(1:显示 2:隐藏)
	Status     *int    `gorm:"column:status;index;"`                     // 状态(1:启用 2:禁用)
	Memo       *string `gorm:"column:status;size:1024;"`                 // 备注
	Creator    *string `gorm:"column:creator;size:36;"`                  // 创建人
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
	item := new(schema.Menu)
	util.StructMapToStruct(a, item)
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

// ----------------------------------------MenuAction--------------------------------------

// GetMenuActionDB 菜单动作
func GetMenuActionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, new(MenuAction))
}

// SchemaMenuAction 菜单动作
type SchemaMenuAction schema.MenuAction

// ToMenuAction 转换为菜单动作实体
func (a SchemaMenuAction) ToMenuAction() *MenuAction {
	item := new(MenuAction)
	util.StructMapToStruct(a, item)
	return item
}

// MenuAction 菜单动作实体
type MenuAction struct {
	Model
	RecordID string  `gorm:"column:record_id;size:36;index;not null;"` // 记录ID
	MenuID   *string `gorm:"column:menu_id;size:36;index;"`            // 菜单ID
	Code     *string `gorm:"column:code;size:100;"`                    // 动作编号
	Name     *string `gorm:"column:name;size:100;"`                    // 动作名称
}

func (a MenuAction) String() string {
	return toString(a)
}

// TableName 表名
func (a MenuAction) TableName() string {
	return a.Model.TableName("menu_action")
}

// ToSchemaMenuAction 转换为菜单动作对象
func (a MenuAction) ToSchemaMenuAction() *schema.MenuAction {
	item := new(schema.MenuAction)
	util.StructMapToStruct(a, item)
	return item
}

// MenuActions 菜单动作列表
type MenuActions []*MenuAction

// ToSchemaMenuActions 转换为菜单动作对象列表
func (a MenuActions) ToSchemaMenuActions() []*schema.MenuAction {
	list := make([]*schema.MenuAction, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuAction()
	}
	return list
}

// ----------------------------------------MenuActionResource--------------------------------------

// GetMenuActionResourceDB 菜单动作关联资源
func GetMenuActionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, new(MenuActionResource))
}

// SchemaMenuActionResource 菜单动作关联资源
type SchemaMenuActionResource schema.MenuActionResource

// ToMenuActionResource 转换为菜单动作关联资源实体
func (a SchemaMenuActionResource) ToMenuActionResource() *MenuActionResource {
	item := new(MenuActionResource)
	util.StructMapToStruct(a, item)
	return item
}

// MenuActionResource 菜单动作关联资源实体
type MenuActionResource struct {
	Model
	RecordID string  `gorm:"column:record_id;size:36;index;not null;"` // 记录ID
	ActionID *string `gorm:"column:action_id;size:36;index;"`          // 菜单动作ID
	Method   *string `gorm:"column:method;size:100;"`                  // 资源请求方式(支持正则)
	Path     *string `gorm:"column:path;size:100;"`                    // 资源请求路径（支持/:id匹配）
}

func (a MenuActionResource) String() string {
	return toString(a)
}

// TableName 表名
func (a MenuActionResource) TableName() string {
	return a.Model.TableName("menu_action_resource")
}

// ToSchemaMenuActionResource 转换为菜单动作关联资源对象
func (a MenuActionResource) ToSchemaMenuActionResource() *schema.MenuActionResource {
	item := new(schema.MenuActionResource)
	util.StructMapToStruct(a, item)
	return item
}

// MenuActionResources 菜单动作关联资源列表
type MenuActionResources []*MenuActionResource

// ToSchemaMenuActionResources 转换为菜单动作关联资源对象列表
func (a MenuActionResources) ToSchemaMenuActionResources() []*schema.MenuActionResource {
	list := make([]*schema.MenuActionResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuActionResource()
	}
	return list
}
