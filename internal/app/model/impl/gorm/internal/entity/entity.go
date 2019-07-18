package entity

import (
	"context"
	"fmt"
	"time"

	icontext "github.com/LyricTian/gin-admin/internal/app/context"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
	"github.com/LyricTian/gin-admin/pkg/util"
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

func getDB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	trans, ok := icontext.FromTrans(ctx)
	if ok {
		db, ok := trans.(*gormplus.DB)
		if ok {
			return db
		}
	}
	return defDB
}

func getDBWithModel(ctx context.Context, defDB *gormplus.DB, m interface{}) *gormplus.DB {
	return gormplus.Wrap(getDB(ctx, defDB).Model(m))
}
