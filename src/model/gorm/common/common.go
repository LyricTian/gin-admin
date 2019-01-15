package common

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
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
	return fmt.Sprintf("%s%s", viper.GetString("db_table_prefix"), name)
}
