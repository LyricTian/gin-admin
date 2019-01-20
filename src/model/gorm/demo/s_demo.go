package demo

import (
	"github.com/LyricTian/gin-admin/src/model/gorm/common"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
)

// GetDemoTableName 获取示例表名
func GetDemoTableName() string {
	return Demo{}.TableName()
}

// Demo 示例程序
type Demo struct {
	common.Model
	RecordID string `gorm:"column:record_id;size:36;unique_index;"` // 记录内码
	Code     string `gorm:"column:code;size:50;index;"`             // 编号
	Name     string `gorm:"column:name;size:100;index;"`            // 名称
	Memo     string `gorm:"column:memo;size:200;"`                  // 备注
	Status   int    `gorm:"column:status;index;"`                   // 状态(1:启用 2:停用)
	Creator  string `gorm:"column:creator;size:36;"`                // 创建者
}

// TableName 表名
func (a Demo) TableName() string {
	return a.Model.TableName("demo")
}

// ToSchemaDemo 转换为模型demo
func (a Demo) ToSchemaDemo() *schema.Demo {
	item := new(schema.Demo)
	_ = util.FillStruct(a, item)
	return item
}

// Demos demo列表
type Demos []*Demo

// ToSchemaDemos 转换为模型demo列表
func (a Demos) ToSchemaDemos() []*schema.Demo {
	var list []*schema.Demo
	_ = util.FillStructs(a, &list)
	return list
}
