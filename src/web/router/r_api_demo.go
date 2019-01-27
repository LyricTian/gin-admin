package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIDemoRouter 注册/demos路由
func APIDemoRouter(g *gin.RouterGroup, c *ctl.Demo) {
	GET(g, "/v1/demos", c.Query, SetName("查询示例数据"))
	GET(g, "/v1/demos/:id", c.Get, SetName("查询指定示例数据"))
	POST(g, "/v1/demos", c.Create, SetName("创建示例数据"))
	PUT(g, "/v1/demos/:id", c.Update, SetName("更新示例数据"))
	DELETE(g, "/v1/demos/:id", c.Delete, SetName("删除示例数据"))
	DELETE(g, "/v1/demos", c.DeleteMany, SetName("删除多条示例数据"))
	PATCH(g, "/v1/demos/:id/enable", c.Enable, SetName("启用示例数据"))
	PATCH(g, "/v1/demos/:id/disable", c.Disable, SetName("禁用示例数据"))
}
