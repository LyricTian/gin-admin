package web

import (
	"context"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/LyricTian/gin-admin/src/web/router"
	"github.com/gin-gonic/gin"
)

// Init 初始化所有服务
func Init(ctx context.Context, obj *inject.Object) *gin.Engine {
	gin.SetMode(config.GetRunMode())
	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	// 注册中间件
	apiPrefixes := []string{"/api/"}
	if dir := config.GetWWWDir(); dir != "" {
		app.Use(middleware.WWWMiddleware(dir, middleware.AllowPathPrefixSkipper(apiPrefixes...)))
	}

	if dir := config.GetSwaggerDir(); dir != "" {
		app.Static("/swagger", dir)
	}

	app.Use(middleware.TraceMiddleware(middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))
	app.Use(middleware.LoggerMiddleware(middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))
	app.Use(middleware.RecoveryMiddleware())

	if config.IsSessionAuth() {
		app.Use(middleware.SessionMiddleware(obj, middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))
	}

	// 注册/api路由
	router.APIHandler(app, obj)

	return app
}
