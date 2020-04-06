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

// Model base model
type Model struct {
	ID        uint       `gorm:"column:id;primary_key;auto_increment;"`
	RecordID  string     `gorm:"column:record_id;size:36;index;default:'';not null;"`
	CreatedAt time.Time  `gorm:"column:created_at;index;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

// TableName table name
func (Model) TableName(name string) string {
	return fmt.Sprintf("%s%s", config.C.Gorm.TablePrefix, name)
}

func toString(v interface{}) string {
	return util.JSONMarshalToString(v)
}

func getDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	trans, ok := icontext.FromTrans(ctx)
	if ok && !icontext.FromNoTrans(ctx) {
		db, ok := trans.(*gorm.DB)
		if ok {
			if icontext.FromTransLock(ctx) {
				if dbType := config.C.Gorm.DBType; dbType == "mysql" ||
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
	db := getDB(ctx, defDB)

	type tabler interface {
		TableName() string
	}

	if t, ok := m.(tabler); ok {
		return db.Table(t.TableName())
	}
	return db.Model(m)
}
