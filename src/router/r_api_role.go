package router

import (
	"github.com/LyricTian/gin-admin/src/api"
	"github.com/LyricTian/gin-admin/src/context"
	"github.com/gin-gonic/gin"
)

// APIRoleRouter 注册/roles路由
func APIRoleRouter(g *gin.RouterGroup, role *api.Role) {
	g.GET("/roles", context.WrapContext(role.Query, "查询角色数据"))
	g.GET("/roles/:id", context.WrapContext(role.Get, "查询指定角色数据"))
	g.POST("/roles", context.WrapContext(role.Create, "创建角色数据"))
	g.PUT("/roles/:id", context.WrapContext(role.Update, "更新角色数据"))
	g.DELETE("/roles/:id", context.WrapContext(role.Delete, "删除角色数据"))
	g.DELETE("/roles", context.WrapContext(role.DeleteMany, "删除多条角色数据"))
	g.PATCH("/roles/:id/enable", context.WrapContext(role.Enable, "启用角色数据"))
	g.PATCH("/roles/:id/disable", context.WrapContext(role.Disable, "禁用角色数据"))
}
