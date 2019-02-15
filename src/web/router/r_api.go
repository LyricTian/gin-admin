package router

import (
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/auth"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/gin-gonic/gin"
)

// APIHandler /api路由
func APIHandler(app *gin.Engine, obj *inject.Object) {
	g := app.Group("/api")

	// 用户授权(session/jwt)
	g.Use(auth.Entry(auth.SkipperFunc(middleware.NoAllowPathPrefixSkipper("/api"))))

	// 用户身份授权
	g.Use(middleware.UserAuthMiddleware(
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
			context.JoinRouter("PUT", "/api/v1/current/password"),
			context.JoinRouter("GET", "/api/v1/current/user"),
			context.JoinRouter("GET", "/api/v1/current/menutree"),
		),
	))

	// 请求频率限制中间件
	if config.GetRateLimiterConfig().Enable {
		g.Use(middleware.RateLimiterMiddleware(obj.RateLimiter))
	}

	c := obj.CtlCommon

	// 注册路由
	APIDemoRouter(g, c.DemoAPI)
	APILoginRouter(g, c.LoginAPI)
	APIRoleRouter(g, c.RoleAPI)
	APIMenuRouter(g, c.MenuAPI)
	APIUserRouter(g, c.UserAPI)
	APIResourceRouter(g, c.ResourceAPI)
}
