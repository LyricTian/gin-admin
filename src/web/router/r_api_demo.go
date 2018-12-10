package router

import (
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIDemoRouter 注册/demos路由
func APIDemoRouter(g *gin.RouterGroup, c *ctl.Demo) {
	g.GET("/demos", context.WrapContext(c.Query, "查询示例数据"))
	g.GET("/demos/:id", context.WrapContext(c.Get, "查询指定示例数据"))
	g.POST("/demos", context.WrapContext(c.Create, "创建示例数据"))
	g.PUT("/demos/:id", context.WrapContext(c.Update, "更新示例数据"))
	g.DELETE("/demos/:id", context.WrapContext(c.Delete, "删除示例数据"))
	g.DELETE("/demos", context.WrapContext(c.DeleteMany, "删除多条示例数据"))
	g.PATCH("/demos/:id/enable", context.WrapContext(c.Enable, "启用示例数据"))
	g.PATCH("/demos/:id/disable", context.WrapContext(c.Disable, "禁用示例数据"))
}
