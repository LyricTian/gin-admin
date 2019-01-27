package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIRoleRouter 注册/roles路由
func APIRoleRouter(g *gin.RouterGroup, c *ctl.Role) {
	GET(g, "/v1/roles", c.Query, SetName("查询角色数据"))
	GET(g, "/v1/roles/:id", c.Get, SetName("查询指定角色数据"))
	POST(g, "/v1/roles", c.Create, SetName("创建角色数据"))
	PUT(g, "/v1/roles/:id", c.Update, SetName("更新角色数据"))
	DELETE(g, "/v1/roles/:id", c.Delete, SetName("删除角色数据"))
	DELETE(g, "/v1/roles", c.DeleteMany, SetName("删除多条角色数据"))
	PATCH(g, "/v1/roles/:id/enable", c.Enable, SetName("启用角色数据"))
	PATCH(g, "/v1/roles/:id/disable", c.Disable, SetName("禁用角色数据"))
}
