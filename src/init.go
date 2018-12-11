package src

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/LyricTian/gin-admin/src/logger"
	model "github.com/LyricTian/gin-admin/src/model/mysql"
	"github.com/LyricTian/gin-admin/src/service/mysql"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/casbin/casbin"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CloseHandle 关闭服务
type CloseHandle func()

// Init 初始化所有服务
func Init(version, traceID string) (*gin.Engine, CloseHandle) {
	// 初始化MySQL
	db := InitMySQL()

	// 初始化日志
	loggerHook := InitLogger(db.Db)

	logger.System(traceID).Infof("服务已运行在[%s]模式下，版本号:%s，进程号：%d",
		viper.GetString("run_mode"), version, os.Getpid())

	// 初始化依赖注入
	enforcer, _, ctlCommon := InitInject(db)

	// 初始化HTTP服务
	httpHandler := web.Init(db, enforcer, ctlCommon)

	return httpHandler, func() {
		// 等待日志钩子写入完成
		if loggerHook != nil {
			loggerHook.Flush()
		}

		// 关闭数据库
		err := db.Close()
		if err != nil {
			logger.System(traceID).Errorf("关闭数据库发生错误: %s", err.Error())
		}
	}
}

// InitInject 初始化依赖注入
func InitInject(db *mysql.DB) (*casbin.Enforcer, *model.Common, *ctl.Common) {
	g := new(inject.Graph)

	// 注入casbin
	enforcer := casbin.NewEnforcer(viper.GetString("casbin_model_conf"), false)
	g.Provide(&inject.Object{Value: enforcer})

	// 注入mysql存储
	modelCommom := new(model.Common).Init(g, db)

	// 注入控制器
	ctlCommon := new(ctl.Common)
	g.Provide(&inject.Object{Value: ctlCommon})

	if err := g.Populate(); err != nil {
		panic("注入模块发生错误:" + err.Error())
	}

	return enforcer, modelCommom, ctlCommon
}

// InitMySQL 初始化mysql数据库
func InitMySQL() *mysql.DB {
	mysqlConfig := viper.GetStringMap("mysql")
	var opts []mysql.Option
	if v := util.T(mysqlConfig["trace"]).Bool(); v {
		opts = append(opts, mysql.SetTrace(v))
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		mysqlConfig["username"],
		mysqlConfig["password"],
		mysqlConfig["addr"],
		mysqlConfig["database"],
	)
	opts = append(opts, mysql.SetDSN(dsn))

	if v := util.T(mysqlConfig["engine"]).String(); v != "" {
		opts = append(opts, mysql.SetEngine(v))
	}

	if v := util.T(mysqlConfig["encoding"]).String(); v != "" {
		opts = append(opts, mysql.SetEncoding(v))
	}

	if v := util.T(mysqlConfig["max_lifetime"]).Int(); v > 0 {
		opts = append(opts, mysql.SetMaxLifetime(time.Duration(v)*time.Second))
	}

	if v := util.T(mysqlConfig["max_open_conns"]).Int(); v > 0 {
		opts = append(opts, mysql.SetMaxOpenConns(v))
	}

	if v := util.T(mysqlConfig["max_idle_conns"]).Int(); v > 0 {
		opts = append(opts, mysql.SetMaxIdleConns(v))
	}

	db, err := mysql.NewDB(opts...)
	if err != nil {
		panic("初始化MySQL数据库发生错误：" + err.Error())
	}

	return db
}

// InitLogger 初始化日志
func InitLogger(mysqlDB *sql.DB) logger.HookFlusher {
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
				mysqlDB,
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
