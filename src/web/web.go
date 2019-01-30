package web

import (
	"context"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/LyricTian/gin-admin/src/web/router"
	"github.com/gin-gonic/gin"
)

// Init 初始化web服务
func Init(ctx context.Context, obj *inject.Object) *gin.Engine {
	span := logger.StartSpan(ctx, "初始化web服务", "web.Init")

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

	// 检查并创建资源数据
	if config.IsAllowCreateResource() {
		err := obj.CtlCommon.CheckAndCreateResource(ctx)
		if err != nil {
			span.Fatalf("检查并创建资源数据发生错误：%s", err.Error())
		}
	}

	// 初始化casbin策略数据
	err := obj.CtlCommon.LoadCasbinPolicyData(ctx)
	if err != nil {
		span.Fatalf("初始化casbin策略数据发生错误：%s", err.Error())
	}

	return app
}
