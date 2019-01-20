package common

import (
	"fmt"
	"time"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/jinzhu/gorm"
)

// Model 定义基础的模型
type Model struct {
	ID        uint       `gorm:"column:id;primary_key;auto_increment;"`
	CreatedAt time.Time  `gorm:"column:created_at;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

// TableName 表名
func (Model) TableName(name string) string {
	return fmt.Sprintf("%s%s", config.GetDBTablePrefix(), name)
}

// WrapPageQuery 包装带有分页的查询
func WrapPageQuery(db *gorm.DB, pp *schema.PaginationParam, out interface{}) (*schema.PaginationResult, error) {
	if pp != nil {
		total, err := gormplus.Wrap(db).FindPage(db, pp.PageIndex, pp.PageSize, out)
		if err != nil {
			return nil, err
		}
		return &schema.PaginationResult{
			Total: total,
		}, nil
	}

	result := db.Find(out)
	return nil, result.Error
}
