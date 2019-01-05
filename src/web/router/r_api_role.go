package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIRoleRouter 注册/roles路由
func APIRoleRouter(g *gin.RouterGroup, c *ctl.Role) {
	GET(g, "/roles", c.Query, "查询角色数据")
	GET(g, "/roles/:id", c.Get, "查询指定角色数据")
	POST(g, "/roles", c.Create, "创建角色数据")
	PUT(g, "/roles/:id", c.Update, "更新角色数据")
	DELETE(g, "/roles/:id", c.Delete, "删除角色数据")
	DELETE(g, "/roles", c.DeleteMany, "删除多条角色数据")
	PATCH(g, "/roles/:id/enable", c.Enable, "启用角色数据")
	PATCH(g, "/roles/:id/disable", c.Disable, "禁用角色数据")
}
