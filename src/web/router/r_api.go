package router

import (
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/gin-gonic/gin"
)

// APIV1Handler /api/v1路由
func APIV1Handler(r *gin.Engine, obj *inject.Object) {
	v1 := r.Group("/api/v1/",
		middleware.VerifySessionMiddleware(
			"/api/v1/login",
			"/api/v1/logout",
		),
		middleware.CasbinMiddleware(obj.Enforcer),
	)

	c := obj.CtlCommon

	// 注册路由
	APIDemoRouter(v1, c.DemoAPI)
	APILoginRouter(v1, c.LoginAPI)
	APIRoleRouter(v1, c.RoleAPI)
	APIMenuRouter(v1, c.MenuAPI)
	APIUserRouter(v1, c.UserAPI)
}
