package initialize

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/pkg/logger"
	loggerhook "github.com/LyricTian/gin-admin/pkg/logger/hook"
	loggergormhook "github.com/LyricTian/gin-admin/pkg/logger/hook/gorm"
	loggermongohook "github.com/LyricTian/gin-admin/pkg/logger/hook/mongo"
)

// InitLogger 初始化日志模块
func InitLogger() (func(), error) {
	c := config.C.Log
	logger.SetLevel(c.Level)
	logger.SetFormatter(c.Format)

	// 设定日志输出
	var file *os.File
	if c.Output != "" {
		switch c.Output {
		case "stdout":
			logger.SetOutput(os.Stdout)
		case "stderr":
			logger.SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
				_ = os.MkdirAll(filepath.Dir(name), 0777)

				f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					return nil, err
				}
				logger.SetOutput(f)
				file = f
			}
		}
	}

	var hook *loggerhook.Hook
	if c.EnableHook {
		switch {
		case c.Hook.IsGorm():
			hc := config.C.LogGormHook

			var dsn string
			switch hc.DBType {
			case "mysql":
				dsn = config.C.MySQL.DSN()
			case "sqlite3":
				dsn = config.C.Sqlite3.DSN()
			case "postgres":
				dsn = config.C.Postgres.DSN()
			default:
				return nil, errors.New("unknown db")
			}

			h := loggerhook.New(loggergormhook.New(&loggergormhook.Config{
				DBType:       hc.DBType,
				DSN:          dsn,
				MaxLifetime:  hc.MaxLifetime,
				MaxOpenConns: hc.MaxOpenConns,
				MaxIdleConns: hc.MaxIdleConns,
				TableName:    hc.Table,
			}),
				loggerhook.SetMaxWorkers(c.HookMaxThread),
				loggerhook.SetMaxQueues(c.HookMaxBuffer),
			)
			logger.AddHook(h)
			hook = h
		case c.Hook.IsMongo():
			h := loggerhook.New(loggermongohook.New(&loggermongohook.Config{
				URI:        config.C.Mongo.URI,
				Database:   config.C.Mongo.Database,
				Timeout:    time.Duration(config.C.Mongo.Timeout) * time.Second,
				Collection: config.C.LogMongoHook.Collection,
			}),
				loggerhook.SetMaxWorkers(c.HookMaxThread),
				loggerhook.SetMaxQueues(c.HookMaxBuffer),
			)
			logger.AddHook(h)
			hook = h
		}
	}

	return func() {
		if file != nil {
			file.Close()
		}

		if hook != nil {
			hook.Flush()
		}
	}, nil
}
