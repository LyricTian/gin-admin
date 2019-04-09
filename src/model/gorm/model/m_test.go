package model_test

import (
	"os"
	"path/filepath"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/spf13/viper"
)

var gdb *gormplus.DB

func init() {
	viper.SetConfigFile("../../../../config/config.toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	cfg := config.GetGorm()
	var dsn string
	switch cfg.DBType {
	case "mysql":
		dsn = config.GetMySQL().DSN()
	case "sqlite3":
		dsn = config.GetSqlite3().DSN()
		os.MkdirAll(filepath.Dir(dsn), 0777)
	case "postgres":
		dsn = config.GetPostgres().DSN()
	default:
		panic("unknown db")
	}

	db, err := gormplus.New(gormplus.Config{
		Debug:        cfg.Debug,
		DBType:       cfg.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.MaxIdleConns,
		MaxLifetime:  cfg.MaxLifetime,
		MaxOpenConns: cfg.MaxOpenConns,
	})
	if err != nil {
		panic(err)
	}
	gdb = db
}
