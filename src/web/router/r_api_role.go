package router

import (
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIRoleRouter 注册/roles路由
func APIRoleRouter(g *gin.RouterGroup, c *ctl.Role) {
	g.GET("/roles", context.WrapContext(c.Query, "查询角色数据"))
	g.GET("/roles/:id", context.WrapContext(c.Get, "查询指定角色数据"))
	g.POST("/roles", context.WrapContext(c.Create, "创建角色数据"))
	g.PUT("/roles/:id", context.WrapContext(c.Update, "更新角色数据"))
	g.DELETE("/roles/:id", context.WrapContext(c.Delete, "删除角色数据"))
	g.DELETE("/roles", context.WrapContext(c.DeleteMany, "删除多条角色数据"))
	g.PATCH("/roles/:id/enable", context.WrapContext(c.Enable, "启用角色数据"))
	g.PATCH("/roles/:id/disable", context.WrapContext(c.Disable, "禁用角色数据"))
}
