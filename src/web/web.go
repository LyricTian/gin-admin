package web

import (
	"context"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/auth"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/LyricTian/gin-admin/src/web/router"
	"github.com/gin-gonic/gin"
	gormstore "github.com/go-session/gorm"
	"github.com/go-session/session"
)

// Init 初始化web服务
func Init(ctx context.Context, obj *inject.Object) *gin.Engine {
	gin.SetMode(config.GetRunMode())

	switch {
	case config.IsSessionAuth():
		auth.SetGlobalAuther(auth.NewSessionAuth(getSessionStore(obj)))
	case config.IsJWTAuth():
		auth.SetGlobalAuther(auth.NewJWTAuth())
	}

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
	app.Use(auth.Entry(auth.SkipperFunc(middleware.NoAllowPathPrefixSkipper(apiPrefixes...))))

	// 注册/api路由
	router.APIHandler(app, obj)

	return app
}

// 获取会话存储
func getSessionStore(obj *inject.Object) session.ManagerStore {
	if config.IsGormDB() && obj.GormDB != nil {
		tableName := config.GetGormTablePrefix() + "session"
		return gormstore.NewStoreWithDB(obj.GormDB.DB, tableName, 0)
	}
	return nil
}
