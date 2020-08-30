package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/v7/internal/app/config"
	"github.com/LyricTian/gin-admin/v7/internal/app/contextx"
	"github.com/jinzhu/gorm"
)

// GetDB ...
func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	trans, ok := contextx.FromTrans(ctx)
	if ok && !contextx.FromNoTrans(ctx) {
		db, ok := trans.(*gorm.DB)
		if ok {
			if contextx.FromTransLock(ctx) {
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

// GetDBWithModel ...
func GetDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return GetDB(ctx, defDB).Model(m)
}
