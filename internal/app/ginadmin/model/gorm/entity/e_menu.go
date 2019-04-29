package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
)

// GetMenuDB 获取菜单存储
func GetMenuDB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	return getDBWithModel(ctx, defDB, Menu{})
}

// GetMenuActionDB 获取菜单动作存储
func GetMenuActionDB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	return getDBWithModel(ctx, defDB, MenuAction{})
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
		Hidden:     &a.Hidden,
		ParentID:   a.ParentID,
		ParentPath: a.ParentPath,
		Creator:    a.Creator,
	}
	return item
}

// ToMenuActions 转换为菜单动作列表
func (a SchemaMenu) ToMenuActions() []*MenuAction {
	list := make([]*MenuAction, len(a.Actions))
	for i, item := range a.Actions {
		list[i] = SchemaMenuAction(*item).ToMenuAction(a.RecordID)
	}
	return list
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
	Hidden     *int   `gorm:"column:hidden;index;"`               // 隐藏菜单(0:不隐藏 1:隐藏)
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
		ParentID:   a.ParentID,
		ParentPath: a.ParentPath,
		Creator:    a.Creator,
		CreatedAt:  &a.CreatedAt,
		UpdatedAt:  &a.UpdatedAt,
	}
	if a.Hidden != nil {
		item.Hidden = *a.Hidden
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

// SchemaMenuAction 菜单动作对象
type SchemaMenuAction schema.MenuAction

// ToMenuAction 转换为菜单动作实体
func (a SchemaMenuAction) ToMenuAction(menuID string) *MenuAction {
	return &MenuAction{
		MenuID: menuID,
		Code:   a.Code,
		Name:   a.Name,
	}
}

// MenuAction 菜单动作关联实体
type MenuAction struct {
	Model
	MenuID string `gorm:"column:menu_id;size:36;index;"` // 菜单ID
	Code   string `gorm:"column:code;size:50;index;"`    // 动作编号
	Name   string `gorm:"column:name;size:50;"`          // 动作名称
}

// TableName 表名
func (a MenuAction) TableName() string {
	return a.Model.TableName("menu_action")
}

// ToSchemaMenuAction 转换为菜单动作对象
func (a MenuAction) ToSchemaMenuAction() *schema.MenuAction {
	return &schema.MenuAction{
		Code: a.Code,
		Name: a.Name,
	}
}

// MenuActions 菜单动作关联实体列表
type MenuActions []*MenuAction

// GetByMenuID 根据菜单ID获取菜单动作列表
func (a MenuActions) GetByMenuID(menuID string) []*schema.MenuAction {
	var list []*schema.MenuAction
	for _, item := range a {
		if item.MenuID == menuID {
			list = append(list, item.ToSchemaMenuAction())
		}
	}
	return list
}

// ToSchemaMenuActions 转换为菜单动作列表
func (a MenuActions) ToSchemaMenuActions() []*schema.MenuAction {
	list := make([]*schema.MenuAction, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuAction()
	}
	return list
}

// ToMap 转换为键值映射
func (a MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

// SchemaMenuResource 菜单资源对象
type SchemaMenuResource schema.MenuResource

// ToMenuResource 转换为菜单资源实体
func (a SchemaMenuResource) ToMenuResource(menuID string) *MenuResource {
	return &MenuResource{
		MenuID: menuID,
		Code:   a.Code,
		Name:   a.Name,
		Method: a.Method,
		Path:   a.Path,
	}
}

// MenuResource 菜单资源关联实体
type MenuResource struct {
	Model
	MenuID string `gorm:"column:menu_id;size:36;index;"` // 菜单ID
	Code   string `gorm:"column:code;size:50;index;"`    // 资源编号
	Name   string `gorm:"column:name;size:50;"`          // 资源名称
	Method string `gorm:"column:method;size:50;"`        // 请求方式
	Path   string `gorm:"column:path;size:255;"`         // 请求路径
}

// TableName 表名
func (a MenuResource) TableName() string {
	return a.Model.TableName("menu_resource")
}

// ToSchemaMenuResource 转换为菜单资源对象
func (a MenuResource) ToSchemaMenuResource() *schema.MenuResource {
	return &schema.MenuResource{
		Code:   a.Code,
		Name:   a.Name,
		Method: a.Method,
		Path:   a.Path,
	}
}

// MenuResources 菜单资源关联实体列表
type MenuResources []*MenuResource

// GetByMenuID 根据菜单ID获取菜单资源列表
func (a MenuResources) GetByMenuID(menuID string) []*schema.MenuResource {
	var list []*schema.MenuResource
	for _, item := range a {
		if item.MenuID == menuID {
			list = append(list, item.ToSchemaMenuResource())
		}
	}
	return list
}

// ToSchemaMenuResources 转换为菜单资源列表
func (a MenuResources) ToSchemaMenuResources() []*schema.MenuResource {
	list := make([]*schema.MenuResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuResource()
	}
	return list
}

// ToMap 转换为键值映射
func (a MenuResources) ToMap() map[string]*MenuResource {
	m := make(map[string]*MenuResource)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}
