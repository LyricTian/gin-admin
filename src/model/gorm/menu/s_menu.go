package gormmenu

import (
	"github.com/LyricTian/gin-admin/src/model/gorm/common"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
)

// GetMenuTableName 获取菜单表名
func GetMenuTableName() string {
	return Menu{}.TableName()
}

// Menu 菜单实体
type Menu struct {
	gormcommon.Model
	RecordID  string `gorm:"column:record_id;size:36;unique_index;"` // 记录内码
	Code      string `gorm:"column:code;size:50;"`                   // 菜单编号
	Name      string `gorm:"column:name;size:50;index;"`             // 菜单名称
	Type      int    `gorm:"column:type;index;"`                     // 菜单类型(1：模块 2：功能 3：资源)
	Sequence  int    `gorm:"column:sequence;index;"`                 // 排序值
	Icon      string `gorm:"column:icon;size:255;"`                  // 菜单图标
	Path      string `gorm:"column:path;size:255;"`                  // 访问路径
	Method    string `gorm:"column:method;size:50;"`                 // 资源请求方式
	LevelCode string `gorm:"column:level_code;size:40;index;"`       // 分级码
	ParentID  string `gorm:"column:parent_id;size:36;index;"`        // 父级内码
	IsHide    int    `gorm:"column:is_hide;index;"`                  // 是否隐藏(1:是 2:否)
	Creator   string `gorm:"column:creator;size:36;"`                // 创建人
}

// TableName 表名
func (a Menu) TableName() string {
	return a.Model.TableName("menu")
}

// ToSchemaMenu 转换为菜单对象
func (a Menu) ToSchemaMenu() *schema.Menu {
	item := new(schema.Menu)
	_ = util.FillStruct(a, item)
	return item
}

// Menus 菜单列表
type Menus []*Menu

// ToSchemaMenus 转换为菜单对象列表
func (a Menus) ToSchemaMenus() []*schema.Menu {
	var list []*schema.Menu
	_ = util.FillStructs(a, &list)
	return list
}
