package entity

import (
	"github.com/LyricTian/gin-admin/src/schema"
)

// GetResourceTableName 获取资源实体表名
func GetResourceTableName() string {
	return Resource{}.TableName()
}

// SchemaResource 资源对象
type SchemaResource schema.Resource

// ToResource 转换为资源实体
func (a SchemaResource) ToResource() *Resource {
	item := &Resource{
		RecordID: a.RecordID,
		Code:     a.Code,
		Name:     a.Name,
		Path:     a.Path,
		Method:   a.Method,
	}
	return item
}

// Resource 资源实体
type Resource struct {
	Model
	RecordID string `gorm:"column:record_id;size:36;unique_index;"` // 记录内码
	Code     string `gorm:"column:code;size:50;index;"`             // 编号
	Name     string `gorm:"column:name;size:100;index;"`            // 名称
	Path     string `gorm:"column:path;size:256;"`                  // 访问路径
	Method   string `gorm:"column:method;size:20;"`                 // 资源请求方式
	Creator  string `gorm:"column:creator;size:36;"`                // 创建者
}

func (a Resource) String() string {
	return toString(a)
}

// TableName 表名
func (a Resource) TableName() string {
	return a.Model.TableName("resource")
}

// ToSchemaResource 转换为资源对象
func (a Resource) ToSchemaResource() *schema.Resource {
	item := &schema.Resource{
		RecordID: a.RecordID,
		Code:     a.Code,
		Name:     a.Name,
		Path:     a.Path,
		Method:   a.Method,
	}
	return item
}

// Resources 资源列表
type Resources []*Resource

// ToSchemaResources 转换为资源对象列表
func (a Resources) ToSchemaResources() []*schema.Resource {
	list := make([]*schema.Resource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaResource()
	}
	return list
}
