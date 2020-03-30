package inject

import (
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/middleware"
	"github.com/LyricTian/gin-admin/internal/app/router"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	// 引入swagger
	_ "github.com/LyricTian/gin-admin/internal/app/swagger"
)

// InitHTTPEngine 初始化gin引擎
func InitHTTPEngine(r *router.Router) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	prefixes := []string{"/api/"}

	// 跟踪ID
	app.Use(middleware.TraceMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// 访问日志
	app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// 崩溃恢复
	app.Use(middleware.RecoveryMiddleware())

	// 跨域请求
	if config.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	// 注册/api路由
	r.RegisterAPI(app)

	// swagger文档
	if config.C.Swagger {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 静态站点
	if dir := config.C.WWW; dir != "" {
		app.Use(middleware.WWWMiddleware(dir))
	}

	return app
}
