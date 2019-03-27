package web

import (
	"sync"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/gin-gonic/gin"
)

var ctxPool = &sync.Pool{
	New: func() interface{} {
		return &context.Context{}
	},
}

func wrapCtx(handler func(ctx *context.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ctxPool.Get().(*context.Context)
		defer ctxPool.Put(ctx)
		ctx.Reset(c)

		handler(ctx)
		c.Abort()
	}
}

// 注册/api路由
func registerAPIRouter(app *gin.Engine, obj *inject.Object) {
	g := app.Group("/api")

	// 用户身份授权
	g.Use(middleware.UserAuthMiddleware(
		obj.Auth,
		middleware.AllowMethodAndPathPrefixSkipper(
			context.JoinRouter("GET", "/api/v1/login"),
			context.JoinRouter("POST", "/api/v1/login"),
		),
	))

	// casbin权限校验中间件
	g.Use(middleware.CasbinMiddleware(obj.Enforcer,
		middleware.AllowMethodAndPathPrefixSkipper(
			context.JoinRouter("GET", "/api/v1/login"),
			context.JoinRouter("POST", "/api/v1/login"),
			context.JoinRouter("POST", "/api/v1/refresh_token"),
			context.JoinRouter("PUT", "/api/v1/current/password"),
			context.JoinRouter("GET", "/api/v1/current/user"),
			context.JoinRouter("GET", "/api/v1/current/menutree"),
		),
	))

	// 请求频率限制中间件
	if config.GetRateLimiter().Enable {
		g.Use(middleware.RateLimiterMiddleware(obj.RateLimiter))
	}

	c := obj.CtlCommon

	// 注册/v1
	v1 := g.Group("/v1")
	{
		// 注册/api/v1/login
		v1.GET("/login/captchaid", wrapCtx(c.LoginCtl.GetCaptchaID))
		v1.GET("/login/captcha", wrapCtx(c.LoginCtl.GetCaptcha))
		v1.POST("/login", wrapCtx(c.LoginCtl.Login))
		v1.POST("/login/exit", wrapCtx(c.LoginCtl.Logout))

		// 注册/api/v1/refresh_token
		v1.POST("/refresh_token", wrapCtx(c.LoginCtl.RefreshToken))

		// 注册/api/v1/current
		v1.PUT("/current/password", wrapCtx(c.LoginCtl.UpdatePassword))
		v1.GET("/current/user", wrapCtx(c.LoginCtl.GetUserInfo))
		v1.GET("/current/menutree", wrapCtx(c.LoginCtl.QueryUserMenuTree))

		// 注册/api/v1/demos
		v1.GET("/demos", wrapCtx(c.DemoCtl.Query))
		v1.GET("/demos/:id", wrapCtx(c.DemoCtl.Get))
		v1.POST("/demos", wrapCtx(c.DemoCtl.Create))
		v1.PUT("/demos/:id", wrapCtx(c.DemoCtl.Update))
		v1.DELETE("/demos/:id", wrapCtx(c.DemoCtl.Delete))
		v1.PATCH("/demos/:id/enable", wrapCtx(c.DemoCtl.Enable))
		v1.PATCH("/demos/:id/disable", wrapCtx(c.DemoCtl.Disable))

		// 注册/api/v1/menus
		v1.GET("/menus", wrapCtx(c.MenuCtl.Query))
		v1.GET("/menus/:id", wrapCtx(c.MenuCtl.Get))
		v1.POST("/menus", wrapCtx(c.MenuCtl.Create))
		v1.PUT("/menus/:id", wrapCtx(c.MenuCtl.Update))
		v1.DELETE("/menus/:id", wrapCtx(c.MenuCtl.Delete))

		// 注册/api/v1/roles
		v1.GET("/roles", wrapCtx(c.RoleCtl.Query))
		v1.GET("/roles/:id", wrapCtx(c.RoleCtl.Get))
		v1.POST("/roles", wrapCtx(c.RoleCtl.Create))
		v1.PUT("/roles/:id", wrapCtx(c.RoleCtl.Update))
		v1.DELETE("/roles/:id", wrapCtx(c.RoleCtl.Delete))

		// 注册/api/v1/users
		v1.GET("/users", wrapCtx(c.UserCtl.Query))
		v1.GET("/users/:id", wrapCtx(c.UserCtl.Get))
		v1.POST("/users", wrapCtx(c.UserCtl.Create))
		v1.PUT("/users/:id", wrapCtx(c.UserCtl.Update))
		v1.DELETE("/users/:id", wrapCtx(c.UserCtl.Delete))
		v1.PATCH("/users/:id/enable", wrapCtx(c.UserCtl.Enable))
		v1.PATCH("/users/:id/disable", wrapCtx(c.UserCtl.Disable))
	}
}
