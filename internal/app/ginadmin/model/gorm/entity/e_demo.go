package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
)

// GetDemoDB 获取demo存储
func GetDemoDB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	return getDBWithModel(ctx, defDB, Demo{})
}

// SchemaDemo demo对象
type SchemaDemo schema.Demo

// ToDemo 转换为demo实体
func (a SchemaDemo) ToDemo() *Demo {
	item := &Demo{
		RecordID: a.RecordID,
		Code:     a.Code,
		Name:     a.Name,
		Memo:     a.Memo,
		Status:   a.Status,
		Creator:  a.Creator,
	}
	return item
}

// Demo demo实体
type Demo struct {
	Model
	RecordID string `gorm:"column:record_id;size:36;index;"` // 记录内码
	Code     string `gorm:"column:code;size:50;index;"`      // 编号
	Name     string `gorm:"column:name;size:100;index;"`     // 名称
	Memo     string `gorm:"column:memo;size:200;"`           // 备注
	Status   int    `gorm:"column:status;index;"`            // 状态(1:启用 2:停用)
	Creator  string `gorm:"column:creator;size:36;"`         // 创建者
}

func (a Demo) String() string {
	return toString(a)
}

// TableName 表名
func (a Demo) TableName() string {
	return a.Model.TableName("demo")
}

// ToSchemaDemo 转换为demo对象
func (a Demo) ToSchemaDemo() *schema.Demo {
	item := &schema.Demo{
		RecordID:  a.RecordID,
		Code:      a.Code,
		Name:      a.Name,
		Memo:      a.Memo,
		Status:    a.Status,
		Creator:   a.Creator,
		CreatedAt: a.CreatedAt,
	}
	return item
}

// Demos demo列表
type Demos []*Demo

// ToSchemaDemos 转换为demo对象列表
func (a Demos) ToSchemaDemos() []*schema.Demo {
	list := make([]*schema.Demo, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaDemo()
	}
	return list
}
