package src

import (
	"context"
	"fmt"
	"os"

	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web"
	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Init 初始化所有服务
func Init(ctx context.Context, version string) (*gin.Engine, func()) {
	// 初始化依赖注入
	obj := inject.Init()

	// 初始化日志
	loggerCallback := InitLogger(obj, version)

	entry := logger.Start(ctx)
	entry.Printf("服务已运行在[%s]模式下，运行版本:%s，进程号：%d",
		viper.GetString("run_mode"), version, os.Getpid())

	// 初始化HTTP服务
	httpHandler := web.Init(obj)
	return httpHandler, func() {
		// 关闭数据库
		if db := obj.MySQL; db != nil {
			if err := db.Close(); err != nil {
				entry.Errorf("关闭数据库发生错误: %s", err.Error())
			}
		}

		// 等待日志钩子写入完成
		if loggerCallback != nil {
			loggerCallback()
		}

	}
}

// InitLogger 初始化日志
func InitLogger(obj *inject.Object, version string) func() {
	logConfig := viper.GetStringMap("log")

	if v := util.T(logConfig["level"]).Int(); v > 0 {
		logrus.SetLevel(logrus.Level(v))
	}

	if v := util.T(logConfig["format"]).String(); v == "json" {
		logrus.SetFormatter(new(logrus.JSONFormatter))
	}

	if v := util.T(logConfig["hook"]).String(); v != "" {
		switch v {
		case "mysql":
			extraItems := []*mysqlhook.ExecExtraItem{
				mysqlhook.NewExecExtraItem(logger.StartTimeKey, "DATETIME"),
				mysqlhook.NewExecExtraItem(logger.UserIDKey, "VARCHAR(36)"),
				mysqlhook.NewExecExtraItem(logger.TraceIDKey, "VARCHAR(100)"),
				mysqlhook.NewExecExtraItem(logger.SpanIDKey, "VARCHAR(100)"),
				mysqlhook.NewExecExtraItem(logger.SpanTitleKey, "VARCHAR(50)"),
				mysqlhook.NewExecExtraItem(logger.SpanFunctionKey, "VARCHAR(200)"),
				mysqlhook.NewExecExtraItem(logger.VersionKey, "VARCHAR(50)"),
			}

			var hookOpts []mysqlhook.Option
			hookOpts = append(hookOpts, mysqlhook.SetExtra(map[string]interface{}{
				logger.VersionKey: version,
			}))

			hookConfig := viper.GetStringMap("log-mysql-hook")
			if v := util.T(hookConfig["max_buffer"]).Int(); v > 0 {
				hookOpts = append(hookOpts, mysqlhook.SetMaxQueues(v))
			}

			if v := util.T(hookConfig["max_thread"]).Int(); v > 0 {
				hookOpts = append(hookOpts, mysqlhook.SetMaxWorkers(v))
			}

			hook := mysqlhook.DefaultWithExtra(
				obj.MySQL.Db,
				fmt.Sprintf("%s_%s",
					viper.GetStringMap("mysql")["table_prefix"],
					util.T(hookConfig["table"]).String()),
				extraItems,
				hookOpts...,
			)

			logrus.AddHook(hook)
			return func() {
				hook.Flush()
			}
		default:
		}
	}

	return nil
}
