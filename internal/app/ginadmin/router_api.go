package ginadmin

import (
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/middleware"
	"github.com/gin-gonic/gin"
)

// 注册/api路由
func registerAPIRouter(app *gin.Engine, obj *Object) {
	g := app.Group("/api")

	// 用户身份授权
	g.Use(middleware.UserAuthMiddleware(
		obj.Auth,
		middleware.AllowMethodAndPathPrefixSkipper(
			middleware.JoinRouter("GET", "/api/v1/login"),
			middleware.JoinRouter("POST", "/api/v1/login"),
		),
	))

	// casbin权限校验中间件
	g.Use(middleware.CasbinMiddleware(obj.Enforcer,
		middleware.AllowMethodAndPathPrefixSkipper(
			middleware.JoinRouter("GET", "/api/v1/login"),
			middleware.JoinRouter("POST", "/api/v1/login"),
			middleware.JoinRouter("POST", "/api/v1/refresh_token"),
			middleware.JoinRouter("PUT", "/api/v1/current/password"),
			middleware.JoinRouter("GET", "/api/v1/current/user"),
			middleware.JoinRouter("GET", "/api/v1/current/menutree"),
		),
	))

	// 请求频率限制中间件
	g.Use(middleware.RateLimiterMiddleware())

	c := obj.Ctl
	v1 := g.Group("/v1")
	{
		// 注册/api/v1/login
		v1.GET("/login/captchaid", c.Login.GetCaptchaID)
		v1.GET("/login/captcha", c.Login.GetCaptcha)
		v1.POST("/login", c.Login.Login)
		v1.POST("/login/exit", c.Login.Logout)

		// 注册/api/v1/refresh_token
		v1.POST("/refresh_token", c.Login.RefreshToken)

		// 注册/api/v1/current
		v1.PUT("/current/password", c.Login.UpdatePassword)
		v1.GET("/current/user", c.Login.GetUserInfo)
		v1.GET("/current/menutree", c.Login.QueryUserMenuTree)

		// 注册/api/v1/demos
		v1.GET("/demos", c.Demo.Query)
		v1.GET("/demos/:id", c.Demo.Get)
		v1.POST("/demos", c.Demo.Create)
		v1.PUT("/demos/:id", c.Demo.Update)
		v1.DELETE("/demos/:id", c.Demo.Delete)
		v1.PATCH("/demos/:id/enable", c.Demo.Enable)
		v1.PATCH("/demos/:id/disable", c.Demo.Disable)

		// 注册/api/v1/menus
		v1.GET("/menus", c.Menu.Query)
		v1.GET("/menus/:id", c.Menu.Get)
		v1.POST("/menus", c.Menu.Create)
		v1.PUT("/menus/:id", c.Menu.Update)
		v1.DELETE("/menus/:id", c.Menu.Delete)

		// 注册/api/v1/roles
		v1.GET("/roles", c.Role.Query)
		v1.GET("/roles/:id", c.Role.Get)
		v1.POST("/roles", c.Role.Create)
		v1.PUT("/roles/:id", c.Role.Update)
		v1.DELETE("/roles/:id", c.Role.Delete)

		// 注册/api/v1/users
		v1.GET("/users", c.User.Query)
		v1.GET("/users/:id", c.User.Get)
		v1.POST("/users", c.User.Create)
		v1.PUT("/users/:id", c.User.Update)
		v1.DELETE("/users/:id", c.User.Delete)
		v1.PATCH("/users/:id/enable", c.User.Enable)
		v1.PATCH("/users/:id/disable", c.User.Disable)
	}
}
