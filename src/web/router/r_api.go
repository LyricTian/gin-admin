package router

import (
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// APIHandler /api路由
func APIHandler(app *gin.Engine, obj *inject.Object) {
	api := app.Group("/api")

	if mode := viper.GetString("auth_mode"); mode == "session" {
		api.Use(middleware.SessionMiddleware(obj))
		api.Use(middleware.VerifySessionMiddleware(
			"/POST/api/v1/login",
			"/POST/api/v1/logout",
		))
	}

	api.Use(middleware.CasbinMiddleware(obj.Enforcer,
		"/POST/api/v1/login",
		"/POST/api/v1/logout",
		"/GET/api/v1/current/menus",
		"/GET/api/v1/current/user",
	))

	c := obj.CtlCommon

	// 注册路由
	APIDemoRouter(api, c.DemoAPI)
	APILoginRouter(api, c.LoginAPI)
	APIRoleRouter(api, c.RoleAPI)
	APIMenuRouter(api, c.MenuAPI)
	APIUserRouter(api, c.UserAPI)
}
