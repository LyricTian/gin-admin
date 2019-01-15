package src

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/service/logrus-hook"
	"github.com/LyricTian/gin-admin/src/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CallbackFunc 回调处理函数
type CallbackFunc func()

// Init 初始化
func Init(ctx context.Context, version string) CallbackFunc {
	// 初始化依赖注入
	obj, err := inject.Init()
	if err != nil {
		panic(err.Error())
	}

	// 初始化日志
	logger.SetVersion(version)
	loggerFunc, err := InitLogger(ctx, obj)
	if err != nil {
		panic(err.Error())
	}

	logger.Start(ctx).Printf("服务开始运行在[%s]模式下，运行版本:%s，进程号：%d",
		viper.GetString("run_mode"), version, os.Getpid())

	// 初始化HTTP服务
	httpFunc := InitHTTPServer(ctx, obj)
	return func() {
		// 等待HTTP服务关闭
		if httpFunc != nil {
			logger.Start(ctx).Printf("关闭HTTP服务")
			httpFunc()
		}

		// 等待日志钩子写入完成
		if loggerFunc != nil {
			logger.Start(ctx).Printf("关闭日志服务")
			loggerFunc()
		}

		// 关闭数据库
		if db := obj.GormDB; db != nil {
			logger.Start(ctx).Printf("关闭数据库服务")
			if err := db.Close(); err != nil {
				logger.Start(ctx).Errorf("关闭数据库发生错误: %s", err.Error())
			}
		}
	}
}

// InitHTTPServer 初始化http服务
func InitHTTPServer(ctx context.Context, obj *inject.Object) CallbackFunc {
	host, port := viper.GetString("http_host"), viper.GetInt("http_port")
	addr := fmt.Sprintf("%s:%d", host, port)

	srv := &http.Server{
		Addr:           addr,
		Handler:        web.Init(obj),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Start(ctx).Printf("HTTP服务开始启动，地址监听在：[%s]", addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Start(ctx).Errorf("监听HTTP服务发生错误: %s", err.Error())
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Start(ctx).Errorf("关闭HTTP服务发生错误: %s", err.Error())
		}
	}
}

// InitLogger 初始化日志
func InitLogger(ctx context.Context, obj *inject.Object) (CallbackFunc, error) {
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

		if v := config.HookMaxThread; v > 0 {
			opts = append(opts, logrushook.SetMaxWorkers(v))
		}
		if v := config.HookMaxBuffer; v > 0 {
			opts = append(opts, logrushook.SetMaxQueues(v))
		}

		if mode := viper.GetString("db_mode"); mode == "gorm" && obj.GormDB != nil {
			hook := logrushook.NewGormHook(obj.GormDB.DB, opts...)
			logrus.AddHook(hook)
			return hook.Flush, nil
		}
	}
	return nil, nil
}
