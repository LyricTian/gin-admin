package web

import (
	"context"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/gin-gonic/gin"
)

// Init 初始化web服务
func Init(ctx context.Context, obj *inject.Object) *gin.Engine {
	gin.SetMode(config.GetRunMode())
	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	apiPrefixes := []string{"/api/"}

	// 跟踪ID
	app.Use(middleware.TraceMiddleware(middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))

	// 访问日志
	app.Use(middleware.LoggerMiddleware(middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))

	// 崩溃恢复
	app.Use(middleware.RecoveryMiddleware())

	// 跨域请求
	if config.GetCORS().Enable {
		app.Use(middleware.CORSMiddleware())
	}

	// 注册/api路由
	registerAPIRouter(app, obj)

	// swagger文档
	if dir := config.GetSwaggerDir(); dir != "" {
		app.Static("/swagger", dir)
	}

	// 静态站点
	if dir := config.GetWWWDir(); dir != "" {
		app.Use(middleware.WWWMiddleware(dir))
	}

	return app
}
