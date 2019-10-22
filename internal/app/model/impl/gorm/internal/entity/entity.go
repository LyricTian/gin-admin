package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/config"
	icontext "github.com/LyricTian/gin-admin/internal/app/context"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/jinzhu/gorm"
)

// 表名前缀
var tablePrefix string

// SetTablePrefix 设定表名前缀
func SetTablePrefix(prefix string) {
	tablePrefix = prefix
}

// GetTablePrefix 获取表名前缀
func GetTablePrefix() string {
	return tablePrefix
}

// Model base model
type Model struct {
	ID        uint       `gorm:"column:id;primary_key;auto_increment;"`
	CreatedAt time.Time  `gorm:"column:created_at;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

// TableName table name
func (Model) TableName(name string) string {
	return fmt.Sprintf("%s%s", GetTablePrefix(), name)
}

func toString(v interface{}) string {
	return util.JSONMarshalToString(v)
}

func getDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	trans, ok := icontext.FromTrans(ctx)
	if ok {
		db, ok := trans.(*gorm.DB)
		if ok {
			if icontext.FromTransLock(ctx) {
				if dbType := config.Global().Gorm.DBType; dbType == "mysql" ||
					dbType == "postgres" {
					db = db.Set("gorm:query_option", "FOR UPDATE")
				}
			}
			return db
		}
	}
	return defDB
}

func getDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return getDB(ctx, defDB).Model(m)
}
