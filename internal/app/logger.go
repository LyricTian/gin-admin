package app

import (
	"os"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v8/internal/app/config"
	"github.com/LyricTian/gin-admin/v8/pkg/logger"
	loggerhook "github.com/LyricTian/gin-admin/v8/pkg/logger/hook"
	loggergormhook "github.com/LyricTian/gin-admin/v8/pkg/logger/hook/gorm"
)

// InitLogger 初始化日志模块
func InitLogger() (func(), error) {
	c := config.C.Log
	logger.SetLevel(logger.Level(c.Level))
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
		var hookLevels []logger.Level
		for _, lvl := range c.HookLevels {
			plvl, err := logger.ParseLevel(lvl)
			if err != nil {
				return nil, err
			}
			hookLevels = append(hookLevels, plvl)
		}

		switch {
		case c.Hook.IsGorm():
			db, err := NewGormDB()
			if err != nil {
				return nil, err
			}

			h := loggerhook.New(loggergormhook.New(db),
				loggerhook.SetMaxWorkers(c.HookMaxThread),
				loggerhook.SetMaxQueues(c.HookMaxBuffer),
				loggerhook.SetLevels(hookLevels...),
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
