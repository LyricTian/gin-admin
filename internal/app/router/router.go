package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v8/internal/app/api"
	"github.com/LyricTian/gin-admin/v8/internal/app/middleware"
	"github.com/LyricTian/gin-admin/v8/pkg/auth"
)

var _ IRouter = (*Router)(nil)

// RouterSet 注入router
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

// IRouter 注册路由
type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

// Router 路由管理器
type Router struct {
	Auth           auth.Auther
	CasbinEnforcer *casbin.SyncedEnforcer
	LoginAPI       *api.LoginAPI
	MenuAPI        *api.MenuAPI
	RoleAPI        *api.RoleAPI
	UserAPI        *api.UserAPI
} // end

// Register 注册路由
func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

// Prefixes 路由前缀列表
func (a *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}

// RegisterAPI register api group router
func (a *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")

	g.Use(middleware.UserAuthMiddleware(a.Auth,
		middleware.AllowPathPrefixSkipper("/api/v1/pub/login"),
	))

	g.Use(middleware.CasbinMiddleware(a.CasbinEnforcer,
		middleware.AllowPathPrefixSkipper("/api/v1/pub"),
	))

	g.Use(middleware.RateLimiterMiddleware())

	v1 := g.Group("/v1")
	{
		pub := v1.Group("/pub")
		{
			gLogin := pub.Group("login")
			{
				gLogin.GET("captchaid", a.LoginAPI.GetCaptcha)
				gLogin.GET("captcha", a.LoginAPI.ResCaptcha)
				gLogin.POST("", a.LoginAPI.Login)
				gLogin.POST("exit", a.LoginAPI.Logout)
			}

			gCurrent := pub.Group("current")
			{
				gCurrent.PUT("password", a.LoginAPI.UpdatePassword)
				gCurrent.GET("user", a.LoginAPI.GetUserInfo)
				gCurrent.GET("menutree", a.LoginAPI.QueryUserMenuTree)
			}
			pub.POST("/refresh-token", a.LoginAPI.RefreshToken)
		}

		gMenu := v1.Group("menus")
		{
			gMenu.GET("", a.MenuAPI.Query)
			gMenu.GET(":id", a.MenuAPI.Get)
			gMenu.POST("", a.MenuAPI.Create)
			gMenu.PUT(":id", a.MenuAPI.Update)
			gMenu.DELETE(":id", a.MenuAPI.Delete)
			gMenu.PATCH(":id/enable", a.MenuAPI.Enable)
			gMenu.PATCH(":id/disable", a.MenuAPI.Disable)
		}
		v1.GET("/menus.tree", a.MenuAPI.QueryTree)

		gRole := v1.Group("roles")
		{
			gRole.GET("", a.RoleAPI.Query)
			gRole.GET(":id", a.RoleAPI.Get)
			gRole.POST("", a.RoleAPI.Create)
			gRole.PUT(":id", a.RoleAPI.Update)
			gRole.DELETE(":id", a.RoleAPI.Delete)
			gRole.PATCH(":id/enable", a.RoleAPI.Enable)
			gRole.PATCH(":id/disable", a.RoleAPI.Disable)
		}
		v1.GET("/roles.select", a.RoleAPI.QuerySelect)

		gUser := v1.Group("users")
		{
			gUser.GET("", a.UserAPI.Query)
			gUser.GET(":id", a.UserAPI.Get)
			gUser.POST("", a.UserAPI.Create)
			gUser.PUT(":id", a.UserAPI.Update)
			gUser.DELETE(":id", a.UserAPI.Delete)
			gUser.PATCH(":id/enable", a.UserAPI.Enable)
			gUser.PATCH(":id/disable", a.UserAPI.Disable)
		}
	} // v1 end
}
