package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIDemoRouter 注册/demos路由
func APIDemoRouter(g *gin.RouterGroup, c *ctl.Demo) {
	GET(g, "/v1/demos", c.Query, SetTitle("查询示例数据"))
	GET(g, "/v1/demos/:id", c.Get, SetTitle("查询指定示例数据"))
	POST(g, "/v1/demos", c.Create, SetTitle("创建示例数据"))
	PUT(g, "/v1/demos/:id", c.Update, SetTitle("更新示例数据"))
	DELETE(g, "/v1/demos/:id", c.Delete, SetTitle("删除示例数据"))
	DELETE(g, "/v1/demos", c.DeleteMany, SetTitle("删除多条示例数据"))
	PATCH(g, "/v1/demos/:id/enable", c.Enable, SetTitle("启用示例数据"))
	PATCH(g, "/v1/demos/:id/disable", c.Disable, SetTitle("禁用示例数据"))
}
