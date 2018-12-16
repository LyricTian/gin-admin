package src

import (
	"fmt"
	"os"

	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web"
	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CloseHandle 关闭服务
type CloseHandle func()

// Init 初始化所有服务
func Init(version, traceID string) (*gin.Engine, CloseHandle) {
	// 初始化依赖注入
	obj := inject.Init()

	// 初始化日志
	loggerHook := InitLogger(obj)

	logger.System(traceID).Infof("服务已运行在[%s]模式下，版本号:%s，进程号：%d",
		viper.GetString("run_mode"), version, os.Getpid())

	// 初始化HTTP服务
	httpHandler := web.Init(obj)

	return httpHandler, func() {
		// 等待日志钩子写入完成
		if loggerHook != nil {
			loggerHook.Flush()
		}

		// 关闭数据库
		if db := obj.MySQL; db != nil {
			if err := db.Close(); err != nil {
				logger.System(traceID).Errorf("关闭数据库发生错误: %s", err.Error())
			}
		}

	}
}

// InitLogger 初始化日志
func InitLogger(obj *inject.Object) logger.HookFlusher {
	logConfig := viper.GetStringMap("log")

	var opts []logger.Option
	if v := util.T(logConfig["level"]).Int(); v > 0 {
		opts = append(opts, logger.SetLevel(v))
	}

	if v := util.T(logConfig["format"]).String(); v != "" {
		opts = append(opts, logger.SetFormat(v))
	}

	l := logger.New(opts...)
	if v := util.T(logConfig["hook"]).String(); v != "" {
		switch v {
		case "mysql":
			extraItems := []*mysqlhook.ExecExtraItem{
				mysqlhook.NewExecExtraItem(logger.FieldKeyType, "varchar(20)"),
				mysqlhook.NewExecExtraItem(logger.FieldKeyUserID, "varchar(36)"),
				mysqlhook.NewExecExtraItem(logger.FieldKeyTraceID, "varchar(100)"),
			}

			var hookOpts []mysqlhook.Option
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

			l.AddHook(hook)
			return hook
		default:
		}
	}

	return nil
}
