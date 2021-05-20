package app

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"os"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v7/internal/app/config"
	"github.com/LyricTian/gin-admin/v7/internal/app/model/gormx"
	"gorm.io/gorm"
)

// InitGormDB 初始化gorm存储
func InitGormDB() (*gorm.DB, func(), error) {
	cfg := config.C.Gorm
	db, cleanFunc, err := NewGormDB()
	if err != nil {
		return nil, cleanFunc, err
	}

	if cfg.EnableAutoMigrate {
		err = gormx.AutoMigrate(db)
		if err != nil {
			return nil, cleanFunc, err
		}
	}

	return db, cleanFunc, nil
}

// NewGormDB 创建DB实例
func NewGormDB() (*gorm.DB, func(), error) {
	cfg := config.C
	var dsn string
	var dialector gorm.Dialector
	switch cfg.Gorm.DBType {
	case "mysql":
		dsn = cfg.MySQL.DSN()
		dialector = mysql.Open(dsn)
	case "sqlite3":
		dsn = cfg.Sqlite3.DSN()
		_ = os.MkdirAll(filepath.Dir(dsn), 0777)
		dialector = sqlite.Open(dsn)
	case "postgres":
		dsn = cfg.Postgres.DSN()
		dialector = postgres.Open(dsn)
	default:
		return nil, nil, errors.New("unknown db")
	}

	return gormx.NewDB(&gormx.Config{
		Debug:        cfg.Gorm.Debug,
		Dialector:    &dialector,
		MaxIdleConns: cfg.Gorm.MaxIdleConns,
		MaxLifetime:  cfg.Gorm.MaxLifetime,
		MaxOpenConns: cfg.Gorm.MaxOpenConns,
		TablePrefix:  cfg.Gorm.TablePrefix,
	})
}
