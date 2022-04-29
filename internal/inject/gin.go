package inject

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/module/consts"
	"github.com/LyricTian/gin-admin/v9/internal/module/ginx"
	"github.com/LyricTian/gin-admin/v9/internal/module/middleware"
	"github.com/LyricTian/gin-admin/v9/internal/router"
	"github.com/LyricTian/gin-admin/v9/internal/schema"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
)

// Initialize gin engine
func InitEngine(r router.IRouter) *gin.Engine {
	// set run mode
	gin.SetMode(config.C.RunMode)

	app := gin.New()

	app.GET("/health", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		ginx.ResOK(c)
	})

	app.NoMethod(func(c *gin.Context) {
		fields := make(map[string]interface{})
		fields["ip"] = c.ClientIP()
		fields["remote_addr"] = c.Request.RemoteAddr
		fields["url"] = c.Request.URL.String()
		fields["proto"] = c.Request.Proto
		fields["user_agent"] = c.GetHeader("User-Agent")

		logger.StandardLogger().WithFields(fields).
			Warnf("No method: %s-%s", c.Request.RequestURI, c.Request.Method)
		ierr, _ := errors.As(errors.MethodNotAllowed(consts.ErrMethodNotAllowedID, ""))
		ginx.ResJSON(c, http.StatusMethodNotAllowed,
			schema.ErrorResult{Error: ierr})
	})

	app.NoRoute(func(c *gin.Context) {
		fields := make(map[string]interface{})
		fields["ip"] = c.ClientIP()
		fields["remote_addr"] = c.Request.RemoteAddr
		fields["url"] = c.Request.URL.String()
		fields["proto"] = c.Request.Proto
		fields["user_agent"] = c.GetHeader("User-Agent")

		logger.StandardLogger().WithFields(fields).
			Warnf("No route: %s-%s", c.Request.RequestURI, c.Request.Method)
		ierr, _ := errors.As(errors.MethodNotAllowed(consts.ErrNotFoundID, ""))
		ginx.ResJSON(c, http.StatusNotFound,
			schema.ErrorResult{Error: ierr})
	})

	prefixes := r.Prefixes()
	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.TraceMiddleware(middleware.AllowedPathPrefix(prefixes...)))
	app.Use(middleware.CopyBodyMiddleware(middleware.AllowedPathPrefix(prefixes...)))
	app.Use(middleware.LoggerMiddleware(middleware.AllowedPathPrefix(prefixes...)))
	app.Use(middleware.CORSMiddleware())
	r.Register(app)

	// Swagger
	if config.C.Swagger {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Website
	if dir := config.C.WWW; dir != "" {
		app.Use(middleware.WWWMiddleware(dir, middleware.DisabledPathPrefix(prefixes...)))
	}

	return app
}
