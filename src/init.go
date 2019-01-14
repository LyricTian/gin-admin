package src

import (
	"context"
	"os"

	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/service/logrus-hook"
	"github.com/LyricTian/gin-admin/src/web"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Init 初始化所有服务
func Init(ctx context.Context, version string) (*gin.Engine, func()) {
	// 初始化依赖注入
	obj, err := inject.Init()
	if err != nil {
		panic(err.Error())
	}

	// 初始化日志
	hookFlush, err := initLogger(obj, version)
	if err != nil {
		panic(err.Error())
	}

	logger.Start(ctx).Printf("服务已运行在[%s]模式下，运行版本:%s，进程号：%d",
		viper.GetString("run_mode"), version, os.Getpid())

	// 初始化HTTP服务
	httpHandler := web.Init(obj)
	return httpHandler, func() {
		// 关闭数据库
		if db := obj.GormDB; db != nil {
			if err := db.Close(); err != nil {
				logger.Start(ctx).Errorf("关闭数据库发生错误: %s", err.Error())
			}
		}

		// 等待日志钩子写入完成
		if hookFlush != nil {
			hookFlush()
		}

	}
}

// 初始化日志
func initLogger(obj *inject.Object, version string) (func(), error) {
	var config struct {
		Level         int    `mapstructure:"level"`
		Format        string `mapstructure:"format"`
		EnableHook    bool   `mapstructure:"enable_hook"`
		HookMaxThread int    `mapstructure:"hook_max_thread"`
		HookMaxBuffer int    `mapstructure:"hook_max_buffer"`
	}

	err := viper.UnmarshalKey("log", &config)
	if err != nil {
		return nil, err
	}

	if v := config.Level; v > -1 {
		logrus.SetLevel(logrus.Level(v))
	}

	if v := config.Format; v == "json" {
		logrus.SetFormatter(new(logrus.JSONFormatter))
	}

	if config.EnableHook {
		var opts []logrushook.Option
		opts = append(opts, logrushook.SetExtra(map[string]interface{}{
			logger.VersionKey: version,
		}))

		if v := config.HookMaxThread; v > 0 {
			opts = append(opts, logrushook.SetMaxWorkers(v))
		}
		if v := config.HookMaxBuffer; v > 0 {
			opts = append(opts, logrushook.SetMaxQueues(v))
		}

		if mode := viper.GetString("db_mode"); mode == "gorm" && obj.GormDB != nil {
			hook := logrushook.NewGormHook(obj.GormDB)
			logrus.AddHook(hook)
			return hook.Flush, nil
		}
	}
	return nil, nil
}
