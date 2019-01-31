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
	span := logger.StartSpanWithCall(ctx, "初始化", "main.src.Init")

	// 依赖注入
	obj, err := inject.Init(ctx)
	if err != nil {
		span().Fatalf("初始化依赖注入发生错误：%s", err.Error())
	}

	// 初始化WEB
	webApp := web.Init(ctx, obj)

	// 初始化数据
	InitData(ctx, obj)

	// 初始化日志钩子
	loggerFunc := InitLoggerHook(ctx, obj)

	// 初始化图形验证码
	if config.IsCaptchaRedisStore() {
		cfg := config.GetRedisConfig()
		captcha.SetCustomStore(store.NewRedisStore(&store.RedisOptions{
			Addr:     cfg.Password,
			Password: cfg.Password,
			DB:       config.GetCaptchaConfig().RedisDB,
		}, captcha.Expiration, log.New(os.Stderr, "[captcha]", log.LstdFlags), "captcha_"))
	}

	// 初始化HTTP服务
	httpFunc := InitHTTPServer(ctx, webApp)
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
func InitHTTPServer(ctx context.Context, handler http.Handler) CallbackFunc {
	span := logger.StartSpanWithCall(ctx, "HTTP服务初始化", "main.src.InitHTTPServer")

	cfg := config.GetHTTPConfig()
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        handler,
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

// InitLoggerHook 初始化日志钩子
func InitLoggerHook(ctx context.Context, obj *inject.Object) CallbackFunc {
	options := config.GetLogConfig()

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

// InitData 初始化数据
func InitData(ctx context.Context, obj *inject.Object) {
	span := logger.StartSpan(ctx, "初始化数据", "main.src.InitData")

	if config.IsAllowCreateResources() {
		// 检查并创建资源数据
		err := obj.CtlCommon.CheckAndCreateResource(ctx)
		if err != nil {
			span.Fatalf("检查并创建资源数据发生错误：%s", err.Error())
		}
	}

	if config.IsAllowInitializeMenus() {
		err := obj.CtlCommon.InitMenuData(ctx)
		if err != nil {
			span.Fatalf("初始化菜单数据发生错误：%s", err.Error())
		}
	}

	// 初始化casbin策略数据
	err := obj.CtlCommon.LoadCasbinPolicyData(ctx)
	if err != nil {
		span.Fatalf("初始化casbin策略数据发生错误：%s", err.Error())
	}
}
