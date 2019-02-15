package web

import (
	"context"
	"log"
	"os"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/captcha/store"
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
	// 初始化图形验证码(redis存储)
	if config.IsCaptchaRedisStore() {
		cfg := config.GetRedisConfig()
		captcha.SetCustomStore(store.NewRedisStore(&store.RedisOptions{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       config.GetCaptchaConfig().RedisDB,
		}, captcha.Expiration,
			log.New(os.Stderr, "[captcha]", log.LstdFlags),
			config.GetCaptchaConfig().RedisPrefix))
	}

	// 初始化认证模式
	switch {
	case config.IsSessionAuth():
		auth.SetGlobalAuther(auth.NewSessionAuth(getSessionStore(obj)))
	case config.IsJWTAuth():
		auth.SetGlobalAuther(auth.NewJWTAuth())
	}

	gin.SetMode(config.GetRunMode())
	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	apiPrefixes := []string{"/api/"}

	// 静态站点
	if dir := config.GetWWWDir(); dir != "" {
		app.Use(middleware.WWWMiddleware(dir, middleware.AllowPathPrefixSkipper(apiPrefixes...)))
	}

	// swagger文档
	if dir := config.GetSwaggerDir(); dir != "" {
		app.Static("/swagger", dir)
	}

	// 跟踪ID
	app.Use(middleware.TraceMiddleware(middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))

	// 访问日志
	app.Use(middleware.LoggerMiddleware(middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))

	// 崩溃恢复
	app.Use(middleware.RecoveryMiddleware())

	// 跨域请求
	if config.GetCORSConfig().Enable {
		app.Use(middleware.CORSMiddleware())
	}

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
