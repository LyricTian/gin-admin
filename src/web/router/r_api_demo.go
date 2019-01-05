package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIDemoRouter 注册/demos路由
func APIDemoRouter(g *gin.RouterGroup, c *ctl.Demo) {
	GET(g, "/demos", c.Query, "查询示例数据")
	GET(g, "/demos/:id", c.Get, "查询指定示例数据")
	POST(g, "/demos", c.Create, "创建示例数据")
	PUT(g, "/demos/:id", c.Update, "更新示例数据")
	DELETE(g, "/demos/:id", c.Delete, "删除示例数据")
	DELETE(g, "/demos", c.DeleteMany, "删除多条示例数据")
	PATCH(g, "/demos/:id/enable", c.Enable, "启用示例数据")
	PATCH(g, "/demos/:id/disable", c.Disable, "禁用示例数据")
}
