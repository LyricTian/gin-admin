package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v9/internal/api"
	"github.com/LyricTian/gin-admin/v9/internal/module/contextx"
	"github.com/LyricTian/gin-admin/v9/internal/module/middleware"
	"github.com/LyricTian/gin-admin/v9/pkg/cache"
	"github.com/LyricTian/gin-admin/v9/pkg/jwtauth"
)

var _ IRouter = (*Router)(nil)

var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

type Router struct {
	Cache          cache.Cacher
	JWTAuth        jwtauth.Auther
	CasbinEnforcer *casbin.Enforcer
	DemoAPI        *api.DemoAPI
} // end

func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

func (a *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}

// RegisterAPI register api group router
func (a *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")

	// Add device_type with request header
	g.Use(func(c *gin.Context) {
		dt := c.Request.Header.Get("DeviceType")
		if dt != "" {
			c.Request = c.Request.WithContext(contextx.NewDeviceType(c.Request.Context(), dt))
		}
		c.Next()
	})

	g.Use(middleware.UserAuthMiddleware(a.JWTAuth, a.Cache,
		middleware.DisabledPathPrefix(
			"/api/v1/verify",
			"/api/v1/login",
		),
	))

	g.Use(middleware.CasbinMiddleware(a.CasbinEnforcer,
		middleware.DisabledPathPrefix(
			"/api/v1/verify",
			"/api/v1/login",
			"/api/v1/current",
		),
	))

	g.Use(middleware.RateLimiterMiddleware())

	v1 := g.Group("/v1")
	{
		gDemo := v1.Group("demos")
		{
			gDemo.GET("", a.DemoAPI.Query)
			gDemo.GET(":id", a.DemoAPI.Get)
			gDemo.POST("", a.DemoAPI.Create)
			gDemo.PUT(":id", a.DemoAPI.Update)
			gDemo.DELETE(":id", a.DemoAPI.Delete)
		}

	} // v1 end
}
