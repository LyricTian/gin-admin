package router

import (
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/gin-gonic/gin"
)

// APIHandler /api路由
func APIHandler(app *gin.Engine, obj *inject.Object) {
	api := app.Group("/api")

	// 用户身份授权
	api.Use(middleware.UserAuthMiddleware(
		middleware.AllowMethodAndPathPrefixSkipper(
			context.JoinRouter("GET", "/api/v1/login"),
			context.JoinRouter("POST", "/api/v1/login"),
		),
	))

	// 权限校验中间件
	api.Use(middleware.CasbinMiddleware(obj.Enforcer,
		middleware.AllowMethodAndPathPrefixSkipper(
			context.JoinRouter("GET", "/api/v1/login"),
			context.JoinRouter("POST", "/api/v1/login"),
			context.JoinRouter("PUT", "/api/v1/current/password"),
			context.JoinRouter("GET", "/api/v1/current/user"),
			context.JoinRouter("GET", "/api/v1/current/menutree"),
		),
	))

	c := obj.CtlCommon

	// 注册路由
	APIDemoRouter(api, c.DemoAPI)
	APILoginRouter(api, c.LoginAPI)
	APIRoleRouter(api, c.RoleAPI)
	APIMenuRouter(api, c.MenuAPI)
	APIUserRouter(api, c.UserAPI)
	APIResourceRouter(api, c.ResourceAPI)
}
