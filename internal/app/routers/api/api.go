package api

import (
	"github.com/LyricTian/gin-admin/internal/app/middleware"
	"github.com/LyricTian/gin-admin/internal/app/routers/api/ctl"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

// RegisterRouter 注册/api路由
func RegisterRouter(app *gin.Engine, container *dig.Container) error {
	err := ctl.Inject(container)
	if err != nil {
		return err
	}

	return container.Invoke(func(
		a auth.Auther,
		e *casbin.SyncedEnforcer,
		cDemo *ctl.Demo,
		cLogin *ctl.Login,
		cMenu *ctl.Menu,
		cRole *ctl.Role,
		cUser *ctl.User,
	) error {

		g := app.Group("/api")

		// 用户身份授权
		g.Use(middleware.UserAuthMiddleware(a,
			middleware.AllowPathPrefixSkipper("/api/v1/pub/login"),
		))

		// casbin权限校验中间件
		g.Use(middleware.CasbinMiddleware(e,
			middleware.AllowPathPrefixSkipper("/api/v1/pub"),
		))

		// 请求频率限制中间件
		g.Use(middleware.RateLimiterMiddleware())

		v1 := g.Group("/v1")
		{
			pub := v1.Group("/pub")
			{
				// 注册/api/v1/pub/login
				gLogin := pub.Group("login")
				{
					gLogin.GET("captchaid", cLogin.GetCaptcha)
					gLogin.GET("captcha", cLogin.ResCaptcha)
					gLogin.POST("", cLogin.Login)
					gLogin.POST("exit", cLogin.Logout)
				}

				// 注册/api/v1/pub/refresh-token
				pub.POST("/refresh-token", cLogin.RefreshToken)

				// 注册/api/v1/pub/current
				gCurrent := pub.Group("current")
				{
					gCurrent.PUT("password", cLogin.UpdatePassword)
					gCurrent.GET("user", cLogin.GetUserInfo)
					gCurrent.GET("menutree", cLogin.QueryUserMenuTree)
				}

			}

			// 注册/api/v1/demos
			gDemo := v1.Group("demos")
			{
				gDemo.GET("", cDemo.Query)
				gDemo.GET(":id", cDemo.Get)
				gDemo.POST("", cDemo.Create)
				gDemo.PUT(":id", cDemo.Update)
				gDemo.DELETE(":id", cDemo.Delete)
				gDemo.PATCH(":id/enable", cDemo.Enable)
				gDemo.PATCH(":id/disable", cDemo.Disable)
			}

			// 注册/api/v1/menus
			gMenu := v1.Group("menus")
			{
				gMenu.GET("", cMenu.Query)
				gMenu.GET(":id", cMenu.Get)
				gMenu.POST("", cMenu.Create)
				gMenu.PUT(":id", cMenu.Update)
				gMenu.DELETE(":id", cMenu.Delete)
			}
			v1.GET("/menus.tree", cMenu.QueryTree)

			// 注册/api/v1/roles
			gRole := v1.Group("roles")
			{
				gRole.GET("", cRole.Query)
				gRole.GET(":id", cRole.Get)
				gRole.POST("", cRole.Create)
				gRole.PUT(":id", cRole.Update)
				gRole.DELETE(":id", cRole.Delete)
			}
			v1.GET("/roles.select", cRole.QuerySelect)

			// 注册/api/v1/users
			gUser := v1.Group("users")
			{
				gUser.GET("", cUser.Query)
				gUser.GET(":id", cUser.Get)
				gUser.POST("", cUser.Create)
				gUser.PUT(":id", cUser.Update)
				gUser.DELETE(":id", cUser.Delete)
				gUser.PATCH(":id/enable", cUser.Enable)
				gUser.PATCH(":id/disable", cUser.Disable)
			}
		}

		return nil
	})
}
