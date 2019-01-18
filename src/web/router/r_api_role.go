package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIRoleRouter 注册/roles路由
func APIRoleRouter(g *gin.RouterGroup, c *ctl.Role) {
	GET(g, "/v1/roles", c.Query, SetTitle("查询角色数据"))
	GET(g, "/v1/roles/:id", c.Get, SetTitle("查询指定角色数据"))
	POST(g, "/v1/roles", c.Create, SetTitle("创建角色数据"))
	PUT(g, "/v1/roles/:id", c.Update, SetTitle("更新角色数据"))
	DELETE(g, "/v1/roles/:id", c.Delete, SetTitle("删除角色数据"))
	DELETE(g, "/v1/roles", c.DeleteMany, SetTitle("删除多条角色数据"))
	PATCH(g, "/v1/roles/:id/enable", c.Enable, SetTitle("启用角色数据"))
	PATCH(g, "/v1/roles/:id/disable", c.Disable, SetTitle("禁用角色数据"))
}
