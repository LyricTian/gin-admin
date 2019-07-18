package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/middleware"
	"github.com/LyricTian/gin-admin/internal/app/routers/api"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

// InitWeb 初始化web引擎
func InitWeb(container *dig.Container) *gin.Engine {
	cfg := config.GetGlobalConfig()
	gin.SetMode(cfg.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	apiPrefixes := []string{"/api/"}

	// 跟踪ID
	app.Use(middleware.TraceMiddleware(middleware.AllowPathPrefixNoSkipper(apiPrefixes...)))

	// 访问日志
	app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(apiPrefixes...)))

	// 崩溃恢复
	app.Use(middleware.RecoveryMiddleware())

	// 跨域请求
	if cfg.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	// 注册/api路由
	err := api.RegisterRouter(app, container)
	handleError(err)

	// swagger文档
	if dir := cfg.Swagger; dir != "" {
		app.Static("/swagger", dir)
	}

	// 静态站点
	if dir := cfg.WWW; dir != "" {
		app.Use(middleware.WWWMiddleware(dir))
	}

	return app
}

// InitHTTPServer 初始化http服务
func InitHTTPServer(ctx context.Context, container *dig.Container) func() {
	cfg := config.GetGlobalConfig().HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      InitWeb(container),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Printf(ctx, "HTTP服务开始启动，地址监听在：[%s]", addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Errorf(ctx, err.Error())
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf(ctx, err.Error())
		}
	}
}
