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

	return container.Invoke(func(a auth.Auther,
		e *casbin.Enforcer,
		demo *ctl.Demo,
		login *ctl.Login,
		menu *ctl.Menu,
		role *ctl.Role,
		user *ctl.User) error {

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
				pub.GET("/login/captchaid", login.GetCaptcha)
				pub.GET("/login/captcha", login.ResCaptcha)
				pub.POST("/login", login.Login)
				pub.POST("/login/exit", login.Logout)

				// 注册/api/v1/pub/refresh_token
				pub.POST("/refresh_token", login.RefreshToken)

				// 注册/api/v1/pub/current
				pub.PUT("/current/password", login.UpdatePassword)
				pub.GET("/current/user", login.GetUserInfo)
				pub.GET("/current/menutree", login.QueryUserMenuTree)
			}

			// 注册/api/v1/demos
			v1.GET("/demos", demo.Query)
			v1.GET("/demos/:id", demo.Get)
			v1.POST("/demos", demo.Create)
			v1.PUT("/demos/:id", demo.Update)
			v1.DELETE("/demos/:id", demo.Delete)
			v1.PATCH("/demos/:id/enable", demo.Enable)
			v1.PATCH("/demos/:id/disable", demo.Disable)

			// 注册/api/v1/menus
			v1.GET("/menus", menu.Query)
			v1.GET("/menus/:id", menu.Get)
			v1.POST("/menus", menu.Create)
			v1.PUT("/menus/:id", menu.Update)
			v1.DELETE("/menus/:id", menu.Delete)

			// 注册/api/v1/roles
			v1.GET("/roles", role.Query)
			v1.GET("/roles/:id", role.Get)
			v1.POST("/roles", role.Create)
			v1.PUT("/roles/:id", role.Update)
			v1.DELETE("/roles/:id", role.Delete)

			// 注册/api/v1/users
			v1.GET("/users", user.Query)
			v1.GET("/users/:id", user.Get)
			v1.POST("/users", user.Create)
			v1.PUT("/users/:id", user.Update)
			v1.DELETE("/users/:id", user.Delete)
			v1.PATCH("/users/:id/enable", user.Enable)
			v1.PATCH("/users/:id/disable", user.Disable)
		}

		return nil
	})
}
