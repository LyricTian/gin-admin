package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

// APIV1Handler /api/v1路由
func APIV1Handler(r *gin.Engine, enforcer *casbin.Enforcer, c *ctl.Common) {
	v1 := r.Group("/api/v1/",
		VerifySessionMiddleware(
			"/api/v1/login",
			"/api/v1/logout",
		),
		CasbinMiddleware(enforcer),
	)

	// 注册路由
	APIDemoRouter(v1, c.DemoAPI)
	APILoginRouter(v1, c.LoginAPI)
	APIRoleRouter(v1, c.RoleAPI)
	APIMenuRouter(v1, c.MenuAPI)
	APIUserRouter(v1, c.UserAPI)
}
