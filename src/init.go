package src

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/captcha/store"
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/service/logrus-hook"
	"github.com/LyricTian/gin-admin/src/web"
	"github.com/sirupsen/logrus"
)

// CallbackFunc 回调处理函数
type CallbackFunc func()

// Init 初始化
func Init(ctx context.Context) CallbackFunc {
	span := logger.StartSpanWithCall(ctx, "服务初始化", "main.Init")

	obj, err := inject.Init(ctx)
	if err != nil {
		span().Fatalf("初始化依赖注入发生错误：%s", err.Error())
	}

	// 初始化日志
	loggerFunc := InitLogger(ctx, obj)

	// 初始化图形验证码
	if config.IsCaptchaRedisStore() {
		config := config.GetRedisConfig()
		captcha.SetCustomStore(store.NewRedisStore(&store.RedisOptions{
			Addr:     config.Password,
			DB:       config.DB,
			Password: config.Password,
		}, captcha.Expiration, log.New(os.Stderr, "[captcha]", log.LstdFlags), "captcha_"))
	}

	// 初始化HTTP服务
	httpFunc := InitHTTPServer(ctx, obj)
	return func() {
		// 等待HTTP服务关闭
		if httpFunc != nil {
			span().Printf("关闭HTTP服务")
			httpFunc()
		}

		// 等待日志钩子写入完成
		if loggerFunc != nil {
			span().Printf("关闭日志服务")
			loggerFunc()
		}

		// 关闭数据库
		if db := obj.GormDB; db != nil {
			span().Printf("关闭数据库服务")
			if err := db.Close(); err != nil {
				span().Errorf("关闭数据库发生错误: %s", err.Error())
			}
		}
	}
}

// InitHTTPServer 初始化http服务
func InitHTTPServer(ctx context.Context, obj *inject.Object) CallbackFunc {
	span := logger.StartSpanWithCall(ctx, "HTTP服务初始化", "main.InitHTTPServer")

	cfg := config.GetHTTPConfig()
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        web.Init(ctx, obj),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		span().Printf("HTTP服务开始启动，地址监听在：[%s]", addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			span().Errorf("监听HTTP服务发生错误: %s", err.Error())
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			span().Errorf("关闭HTTP服务发生错误: %s", err.Error())
		}
	}
}

// InitLogger 初始化日志
func InitLogger(ctx context.Context, obj *inject.Object) CallbackFunc {
	options := config.GetLogConfig()
	if v := options.Level; v > -1 {
		logrus.SetLevel(logrus.Level(v))
	}

	if v := options.Format; v == "json" {
		logrus.SetFormatter(new(logrus.JSONFormatter))
	}

	if options.EnableHook {
		var opts []logrushook.Option

		if v := options.HookMaxThread; v > 0 {
			opts = append(opts, logrushook.SetMaxWorkers(v))
		}
		if v := options.HookMaxBuffer; v > 0 {
			opts = append(opts, logrushook.SetMaxQueues(v))
		}

		if config.IsGormDB() && obj.GormDB != nil {
			hook := logrushook.NewGormHook(obj.GormDB.DB, opts...)
			logrus.AddHook(hook)
			return hook.Flush
		}
	}
	return nil
}
