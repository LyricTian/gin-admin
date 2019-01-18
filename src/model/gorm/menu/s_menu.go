package menu

import (
	"github.com/LyricTian/gin-admin/src/model/gorm/common"
)

// GetMenuTableName 获取菜单表名
func GetMenuTableName() string {
	return Menu{}.TableName()
}

// Menu 菜单管理
type Menu struct {
	common.Model
	RecordID  string `gorm:"column:record_id;size:36;unique_index;"` // 记录内码
	Code      string `gorm:"column:code;size:50;"`                   // 菜单编号
	Name      string `gorm:"column:name;size:50;index;"`             // 菜单名称
	Type      int    `gorm:"column:type;index;"`                     // 菜单类型(10：模块 20：功能 30：资源)
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
