package main

import (
	"database/sql"
	"flag"
	"fmt"
	"gin-admin/src/api"
	"gin-admin/src/context"
	"gin-admin/src/logger"
	model "gin-admin/src/model/mysql"
	"gin-admin/src/router"
	"gin-admin/src/service/mysql"
	"gin-admin/src/util"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

// VERSION 当前服务版本号
const VERSION = "1.0.0"

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "配置文件(.json,.yaml,.toml)")
	flag.StringVar(&configFile, "c", "", "配置文件(.json,.yaml,.toml)")
}

func main() {
	flag.Parse()

	if configFile == "" {
		panic("请使用-c指定配置文件")
	}

	// 初始化配置文件
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件发生错误：" + err.Error())
	}

	traceID := util.UUIDString()

	// 初始化MySQL数据库
	mysqlDB := initMySQL()

	// 初始化日志
	loggerHook := initLogger(mysqlDB.Db)

	logger.System(traceID).Infof("服务已运行在[%s]模式下，版本号:%s，进程号：%d", viper.GetString("run_mode"), VERSION, os.Getpid())

	g := new(inject.Graph)

	// 注入mysql存储
	new(model.Common).Init(g, mysqlDB)

	// 注入API
	apiCommon := new(api.Common)
	g.Provide(&inject.Object{Value: apiCommon})

	if err := g.Populate(); err != nil {
		logger.System(traceID).Panicf("注入模块发生错误:%v", err)
	}

	var state int32 = 1
	ac := make(chan error)
	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGTERM, syscall.SIGQUIT)

	// 开启HTTP监听
	httpServer := initHTTPServer(apiCommon)

	go func() {
		logger.System(traceID).Infof("HTTP服务启动成功，端口监听在[%s]", viper.GetString("http_addr"))
		ac <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-ac:
		if err != nil && atomic.LoadInt32(&state) == 1 {
			logger.System(traceID).Errorf("监听HTTP服务发生错误:%s", err.Error())
		}
	case sig := <-sc:
		atomic.StoreInt32(&state, 0)
		logger.System(traceID).Infof("获取到退出信号[%s]", sig.String())
	}

	// 等待日志钩子写入完成
	if loggerHook != nil {
		loggerHook.Flush()
	}

	// 关闭MySQL数据库
	if err := mysqlDB.Close(); err != nil {
		logger.System(traceID).Errorf("关闭数据库发生错误:%s", err.Error())
	}

	// 退出应用
	os.Exit(int(atomic.LoadInt32(&state)))
}

// 初始化HTTP服务
func initHTTPServer(apiCommon *api.Common) *http.Server {
	gin.SetMode(viper.GetString("run_mode"))

	app := gin.New()

	// 注册中间件
	app.Use(context.WrapContext(router.TraceMiddleware))
	app.Use(logger.Middleware("/api/"))
	app.Use(context.WrapContext(router.RecoveryMiddleware))

	app.NoMethod(context.WrapContext(func(ctx *context.Context) {
		ctx.ResError(fmt.Errorf("方法不允许"), 405)
	}))

	app.NoRoute(context.WrapContext(func(ctx *context.Context) {
		ctx.ResError(fmt.Errorf("资源不存在"), 404)
	}))

	// 注册/api/v1路由
	router.APIV1Handler(app, apiCommon)

	return &http.Server{
		Addr:           viper.GetString("http_addr"),
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// 初始化mysql数据库
func initMySQL() *mysql.DB {
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

// 初始化日志
func initLogger(mysqlDB *sql.DB) logger.HookFlusher {
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
