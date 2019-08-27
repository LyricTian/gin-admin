package api

import (
	"github.com/LyricTian/gin-admin/internal/app/middleware"
	"github.com/LyricTian/gin-admin/internal/app/routers/api/ctl"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/casbin/casbin"
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
		e *casbin.Enforcer,
		cDemo *ctl.Demo,
		cLogin *ctl.Login,
		cMenu *ctl.Menu,
		cRole *ctl.Role,
		cUser *ctl.User,
	) error {

		g := app.Group("/api")

		// 用户身份授权
		g.Use(middleware.UserAuthMiddleware(
			a,
			middleware.AllowMethodAndPathPrefixSkipper(
				middleware.JoinRouter("GET", "/api/v1/pub/login"),
				middleware.JoinRouter("POST", "/api/v1/pub/login"),
			),
		))

		// casbin权限校验中间件
		g.Use(middleware.CasbinMiddleware(e,
			middleware.AllowMethodAndPathPrefixSkipper(
				middleware.JoinRouter("GET", "/api/v1/pub"),
				middleware.JoinRouter("POST", "/api/v1/pub"),
			),
		))

		// 请求频率限制中间件
		g.Use(middleware.RateLimiterMiddleware())

		v1 := g.Group("/v1")
		{
			pub := v1.Group("/pub")
			{
				// 注册/api/v1/pub/login
				pub.GET("/login/captchaid", cLogin.GetCaptcha)
				pub.GET("/login/captcha", cLogin.ResCaptcha)
				pub.POST("/login", cLogin.Login)
				pub.POST("/login/exit", cLogin.Logout)

				// 注册/api/v1/pub/refresh_token
				pub.POST("/refresh_token", cLogin.RefreshToken)

				// 注册/api/v1/pub/current
				pub.PUT("/current/password", cLogin.UpdatePassword)
				pub.GET("/current/user", cLogin.GetUserInfo)
				pub.GET("/current/menutree", cLogin.QueryUserMenuTree)
			}

			// 注册/api/v1/demos
			v1.GET("/demos", cDemo.Query)
			v1.GET("/demos/:id", cDemo.Get)
			v1.POST("/demos", cDemo.Create)
			v1.PUT("/demos/:id", cDemo.Update)
			v1.DELETE("/demos/:id", cDemo.Delete)
			v1.PATCH("/demos/:id/enable", cDemo.Enable)
			v1.PATCH("/demos/:id/disable", cDemo.Disable)

			// 注册/api/v1/menus
			v1.GET("/menus", cMenu.Query)
			v1.GET("/menus/:id", cMenu.Get)
			v1.POST("/menus", cMenu.Create)
			v1.PUT("/menus/:id", cMenu.Update)
			v1.DELETE("/menus/:id", cMenu.Delete)

			// 注册/api/v1/roles
			v1.GET("/roles", cRole.Query)
			v1.GET("/roles/:id", cRole.Get)
			v1.POST("/roles", cRole.Create)
			v1.PUT("/roles/:id", cRole.Update)
			v1.DELETE("/roles/:id", cRole.Delete)

			// 注册/api/v1/users
			v1.GET("/users", cUser.Query)
			v1.GET("/users/:id", cUser.Get)
			v1.POST("/users", cUser.Create)
			v1.PUT("/users/:id", cUser.Update)
			v1.DELETE("/users/:id", cUser.Delete)
			v1.PATCH("/users/:id/enable", cUser.Enable)
			v1.PATCH("/users/:id/disable", cUser.Disable)
		}

		return nil
	})
}
