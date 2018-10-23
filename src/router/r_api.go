package router

import (
	"gin-admin/src/api"

	"github.com/gin-gonic/gin"
)

// APIV1Handler /api/v1路由
func APIV1Handler(r *gin.Engine, c *api.Common) {
	v1 := r.Group("/api/v1")
	APIDemoRouter(v1, c.DemoAPI)
	APIMenuRouter(v1, c.MenuAPI)
}
