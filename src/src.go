package src

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/gops/agent"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/captcha/store"
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/web"
)

// ReleaseFunc 资源释放函数
type ReleaseFunc func()

// Init 初始化
func Init(ctx context.Context) ReleaseFunc {
	span := logger.StartSpanWithCall(ctx, "服务初始化", "main.src.Init")

	// 依赖注入
	obj, err := inject.Init(ctx)
	if err != nil {
		span().Fatalf("依赖注入初始化发生错误：%s", err.Error())
	}

	// 初始化菜单数据
	if config.GetAllowInitMenu() {
		if err := obj.CtlCommon.InitMenuData(ctx); err != nil {
			span().Fatalf("初始化菜单数据发生错误：%s", err.Error())
		}
	}

	// 加载casbin策略数据
	if err := obj.CtlCommon.LoadCasbinPolicyData(ctx); err != nil {
		span().Fatalf("加载casbin策略数据发生错误：%s", err.Error())
	}

	// 图形验证码(redis存储)
	if c := config.GetCaptcha(); c.Store == "redis" {
		rc := config.GetRedis()
		captcha.SetCustomStore(store.NewRedisStore(&store.RedisOptions{
			Addr:     rc.Addr,
			Password: rc.Password,
			DB:       c.RedisDB,
		}, captcha.Expiration, logger.StandardLogger(), c.RedisPrefix))
	}

	// HTTP服务
	httpRFunc := httpServerInit(ctx, obj)

	// 服务监控
	if c := config.GetMonitor(); c.Enable {
		err = agent.Listen(agent.Options{Addr: c.Addr, ConfigDir: c.ConfigDir})
		if err != nil {
			span().Errorf("开启监控发生错误：%s", err.Error())
		}
	}

	return func() {
		// 等待HTTP服务关闭
		if httpRFunc != nil {
			span().Printf("关闭HTTP服务")
			httpRFunc()
		}

		if a := obj.Auth; a != nil {
			if err := a.Release(); err != nil {
				span().Errorf("释放认证资源发生错误: %s", err.Error())
			}
		}

		// 关闭数据库
		if db := obj.GormDB; db != nil {
			span().Printf("关闭数据库服务")
			if err := db.Close(); err != nil {
				span().Errorf("关闭数据库发生错误: %s", err.Error())
			}
		}

		// 如果启用监控，则关闭
		if config.GetMonitor().Enable {
			agent.Close()
		}
	}
}

// httpServerInit HTTP服务初始化
func httpServerInit(ctx context.Context, obj *inject.Object) ReleaseFunc {
	span := logger.StartSpanWithCall(ctx, "HTTP服务初始化", "main.src.httpServerInit")

	cfg := config.GetHTTP()
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      web.Init(ctx, obj),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
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

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			span().Errorf("关闭HTTP服务发生错误: %s", err.Error())
		}
	}
}
