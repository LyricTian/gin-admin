package router

import (
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/gin-gonic/gin"
)

// APIHandler /api路由
func APIHandler(app *gin.Engine, obj *inject.Object) {
	api := app.Group("/api")

	switch {
	case config.IsSessionAuth():
		api.Use(middleware.VerifySessionMiddleware(
			middleware.AllowMethodAndPathPrefixSkipper(
				context.GetRouter("GET", "/api/v1/login"),
				context.GetRouter("POST", "/api/v1/login"),
			),
		))
	}

	// 权限校验中间件
	api.Use(middleware.CasbinMiddleware(obj.Enforcer,
		middleware.AllowMethodAndPathPrefixSkipper(
			context.GetRouter("GET", "/api/v1/login"),
			context.GetRouter("POST", "/api/v1/login"),
			context.GetRouter("PUT", "/api/v1/current/password"),
			context.GetRouter("GET", "/api/v1/current/user"),
			context.GetRouter("GET", "/api/v1/current/menutree"),
		),
	))

	c := obj.CtlCommon

	// 注册路由
	APIDemoRouter(api, c.DemoAPI)
	APILoginRouter(api, c.LoginAPI)
	APIRoleRouter(api, c.RoleAPI)
	APIMenuRouter(api, c.MenuAPI)
	APIUserRouter(api, c.UserAPI)
}
