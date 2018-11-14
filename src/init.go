package src

import (
	"database/sql"
	"fmt"
	"gin-admin/src/api"
	"gin-admin/src/context"
	"gin-admin/src/logger"
	model "gin-admin/src/model/mysql"
	"gin-admin/src/router"
	"gin-admin/src/service/mysql"
	"gin-admin/src/util"
	"os"
	"time"

	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CloseHandle 关闭服务
type CloseHandle func()

// Init 初始化所有服务
func Init(version, traceID string) (*gin.Engine, CloseHandle) {
	db := InitMySQL()

	loggerHook := InitLogger(db.Db)

	logger.System(traceID).Infof("服务已运行在[%s]模式下，版本号:%s，进程号：%d",
		viper.GetString("run_mode"), version, os.Getpid())

	apiCommon := InitInject(db)
	httpHandler := InitHTTPHandler(apiCommon, db)

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
func InitInject(db *mysql.DB) *api.Common {
	g := new(inject.Graph)

	// 注入mysql存储
	new(model.Common).Init(g, db)

	// 注入API
	apiCommon := new(api.Common)
	g.Provide(&inject.Object{Value: apiCommon})

	if err := g.Populate(); err != nil {
		panic("注入模块发生错误:" + err.Error())
	}

	return apiCommon
}

// InitHTTPHandler 初始化GIN服务
func InitHTTPHandler(apiCommon *api.Common, db *mysql.DB) *gin.Engine {
	gin.SetMode(viper.GetString("run_mode"))

	app := gin.New()

	// 注册中间件
	apiPrefixes := []string{
		"/api/",
	}

	app.Use(router.TraceMiddleware(apiPrefixes...))
	app.Use(logger.Middleware(apiPrefixes...))
	app.Use(router.RecoveryMiddleware)
	app.Use(router.SessionMiddleware(db, apiPrefixes...))

	app.NoMethod(context.WrapContext(func(ctx *context.Context) {
		ctx.ResError(fmt.Errorf("方法不允许"), 405)
	}))

	app.NoRoute(context.WrapContext(func(ctx *context.Context) {
		ctx.ResError(fmt.Errorf("资源不存在"), 404)
	}))

	// 注册/api/v1路由
	router.APIV1Handler(app, apiCommon)

	return app
}

// InitMySQL 初始化mysql数据库
func InitMySQL() *mysql.DB {
	mysqlConfig := viper.GetStringMap("mysql")
	var opts []mysql.Option
	if v := util.T(mysqlConfig["trace"]).Bool(); v {
		opts = append(opts, mysql.SetTrace(v))
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		mysqlConfig["username"], mysqlConfig["password"], mysqlConfig["addr"], mysqlConfig["database"],
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
				mysqlhook.NewExecExtraItem(logger.FieldKeyTraceID, "varchar(36)"),
			}

			var hookOpts []mysqlhook.Option
			hookConfig := viper.GetStringMap("log-mysql-hook")
			if v := util.T(hookConfig["max_buffer"]).Int(); v > 0 {
				hookOpts = append(hookOpts, mysqlhook.SetMaxQueues(v))
			}

			if v := util.T(hookConfig["max_thread"]).Int(); v > 0 {
				hookOpts = append(hookOpts, mysqlhook.SetMaxWorkers(v))
			}

			hook := mysqlhook.DefaultWithExtra(mysqlDB, fmt.Sprintf("%s_%s", viper.GetString("mysql_table_prefix"), util.T(hookConfig["table"]).String()), extraItems, hookOpts...)
			l.AddHook(hook)
			return hook
		default:
		}
	}

	return nil
}
