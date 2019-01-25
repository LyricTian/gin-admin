package inject

import (
	"errors"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
)

// 获取gorm存储
func getGormDB() (*gormplus.DB, error) {
	cfg := config.GetGormConfig()

	var dsn string
	switch cfg.DBType {
	case "mysql":
		dsn = config.GetMySQLConfig().DSN()
	case "sqlite3":
		dsn = config.GetSqlite3Config().DSN()
	case "postgres":
		dsn = config.GetPostgresConfig().DSN()
	default:
		return nil, errors.New("unknown db")
	}

	return gormplus.New(gormplus.Config{
		Debug:        cfg.Debug,
		DBType:       cfg.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.MaxIdleConns,
		MaxLifetime:  cfg.MaxLifetime,
		MaxOpenConns: cfg.MaxOpenConns,
	})
}
