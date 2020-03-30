package inject

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/LyricTian/gin-admin/internal/app/config"
	igorm "github.com/LyricTian/gin-admin/internal/app/model/impl/gorm"
	"github.com/jinzhu/gorm"
)

// InitGormDB 初始化存储
func InitGormDB() (*gorm.DB, func(), error) {
	cfg := config.C.Gorm
	db, err := getGormDB()
	if err != nil {
		return nil, nil, err
	}

	cleanFunc := func() {
		db.Close()
	}
	igorm.SetTablePrefix(cfg.TablePrefix)
	if cfg.EnableAutoMigrate {
		err = igorm.AutoMigrate(db)
		if err != nil {
			return nil, cleanFunc, err
		}
	}

	return db, cleanFunc, nil
}

// getGormDB 获取gorm存储
func getGormDB() (*gorm.DB, error) {
	cfg := config.C

	var dsn string
	switch cfg.Gorm.DBType {
	case "mysql":
		dsn = cfg.MySQL.DSN()
	case "sqlite3":
		dsn = cfg.Sqlite3.DSN()
		_ = os.MkdirAll(filepath.Dir(dsn), 0777)
	case "postgres":
		dsn = cfg.Postgres.DSN()
	default:
		return nil, errors.New("unknown db")
	}

	return igorm.NewDB(&igorm.Config{
		Debug:        cfg.Gorm.Debug,
		DBType:       cfg.Gorm.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.Gorm.MaxIdleConns,
		MaxLifetime:  cfg.Gorm.MaxLifetime,
		MaxOpenConns: cfg.Gorm.MaxOpenConns,
	})
}
