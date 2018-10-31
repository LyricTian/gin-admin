package router

import (
	"gin-admin/src/api"

	"github.com/gin-gonic/gin"
)

// APIV1Handler /api/v1路由
func APIV1Handler(r *gin.Engine, c *api.Common) {
	relativePath := "/api/v1"
	v1 := r.Group(relativePath,
		VerifySessionMiddleware(
			[]string{relativePath + "/"},
			[]string{
				relativePath + "/login",
				relativePath + "/logout",
			},
		))

	// 注册路由
	APILoginRouter(v1, c.LoginAPI)
	APIRoleRouter(v1, c.RoleAPI)
	APIDemoRouter(v1, c.DemoAPI)
	APIMenuRouter(v1, c.MenuAPI)
	APIUserRouter(v1, c.UserAPI)
}
