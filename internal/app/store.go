package app

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/LyricTian/gin-admin/internal/app/config"
	igorm "github.com/LyricTian/gin-admin/internal/app/model/impl/gorm"
	"github.com/jinzhu/gorm"
	"go.uber.org/dig"
)

// InitStore 初始化存储
func InitStore(container *dig.Container) (func(), error) {
	var storeCall func()
	cfg := config.Global()

	switch cfg.Store {
	case "gorm":
		db, err := initGorm()
		if err != nil {
			return nil, err
		}

		storeCall = func() {
			db.Close()
		}

		igorm.SetTablePrefix(cfg.Gorm.TablePrefix)

		if cfg.Gorm.EnableAutoMigrate {
			err = igorm.AutoMigrate(db)
			if err != nil {
				return nil, err
			}
		}

		// 注入DB
		container.Provide(func() *gorm.DB {
			return db
		})

		igorm.Inject(container)
	default:
		return nil, errors.New("unknown store")
	}

	return storeCall, nil
}

// initGorm 实例化gorm存储
func initGorm() (*gorm.DB, error) {
	cfg := config.Global()

	var dsn string
	switch cfg.Gorm.DBType {
	case "mysql":
		dsn = cfg.MySQL.DSN()
	case "sqlite3":
		dsn = cfg.Sqlite3.DSN()
		os.MkdirAll(filepath.Dir(dsn), 0777)
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
