package api

import (
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/bll"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/middleware"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/routers/api/ctl"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

// RegisterRouter 注册/api路由
func RegisterRouter(app *gin.Engine, b *bll.Common, a auth.Auther, enforcer *casbin.Enforcer) {
	g := app.Group("/api")

	// 用户身份授权
	g.Use(middleware.UserAuthMiddleware(
		a,
		middleware.AllowMethodAndPathPrefixSkipper(
			middleware.JoinRouter("GET", "/api/v1/login"),
			middleware.JoinRouter("POST", "/api/v1/login"),
		),
	))

	// casbin权限校验中间件
	g.Use(middleware.CasbinMiddleware(enforcer,
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

	demoCtl := ctl.NewDemo(b)
	loginCtl := ctl.NewLogin(b)
	menuCtl := ctl.NewMenu(b)
	roleCtl := ctl.NewRole(b)
	userCtl := ctl.NewUser(b)

	v1 := g.Group("/v1")
	{
		// 注册/api/v1/login
		v1.GET("/login/captchaid", loginCtl.GetCaptchaID)
		v1.GET("/login/captcha", loginCtl.GetCaptcha)
		v1.POST("/login", loginCtl.Login)
		v1.POST("/login/exit", loginCtl.Logout)

		// 注册/api/v1/refresh_token
		v1.POST("/refresh_token", loginCtl.RefreshToken)

		// 注册/api/v1/current
		v1.PUT("/current/password", loginCtl.UpdatePassword)
		v1.GET("/current/user", loginCtl.GetUserInfo)
		v1.GET("/current/menutree", loginCtl.QueryUserMenuTree)

		// 注册/api/v1/demos
		v1.GET("/demos", demoCtl.Query)
		v1.GET("/demos/:id", demoCtl.Get)
		v1.POST("/demos", demoCtl.Create)
		v1.PUT("/demos/:id", demoCtl.Update)
		v1.DELETE("/demos/:id", demoCtl.Delete)
		v1.PATCH("/demos/:id/enable", demoCtl.Enable)
		v1.PATCH("/demos/:id/disable", demoCtl.Disable)

		// 注册/api/v1/menus
		v1.GET("/menus", menuCtl.Query)
		v1.GET("/menus/:id", menuCtl.Get)
		v1.POST("/menus", menuCtl.Create)
		v1.PUT("/menus/:id", menuCtl.Update)
		v1.DELETE("/menus/:id", menuCtl.Delete)

		// 注册/api/v1/roles
		v1.GET("/roles", roleCtl.Query)
		v1.GET("/roles/:id", roleCtl.Get)
		v1.POST("/roles", roleCtl.Create)
		v1.PUT("/roles/:id", roleCtl.Update)
		v1.DELETE("/roles/:id", roleCtl.Delete)

		// 注册/api/v1/users
		v1.GET("/users", userCtl.Query)
		v1.GET("/users/:id", userCtl.Get)
		v1.POST("/users", userCtl.Create)
		v1.PUT("/users/:id", userCtl.Update)
		v1.DELETE("/users/:id", userCtl.Delete)
		v1.PATCH("/users/:id/enable", userCtl.Enable)
		v1.PATCH("/users/:id/disable", userCtl.Disable)
	}
}
