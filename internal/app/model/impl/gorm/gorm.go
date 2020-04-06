package gorm

import (
	"context"
	"strings"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/entity"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/jinzhu/gorm"

	// gorm存储注入
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Config 配置参数
type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

// NewDB 创建DB实例
func NewDB(c *Config) (*gorm.DB, func(), error) {
	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		return nil, nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	cleanFunc := func() {
		err := db.Close()
		if err != nil {
			logger.Errorf(context.Background(), "Gorm db close error: %s", err.Error())
		}
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, cleanFunc, err
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return db, cleanFunc, nil
}

// AutoMigrate 自动映射数据表
func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(entity.Demo),
		new(entity.MenuAction),
		new(entity.MenuActionResource),
		new(entity.Menu),
		new(entity.RoleMenu),
		new(entity.Role),
		new(entity.UserRole),
		new(entity.User),
	).Error
}
