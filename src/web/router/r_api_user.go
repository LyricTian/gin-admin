package router

import (
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIUserRouter 注册/users路由
func APIUserRouter(g *gin.RouterGroup, c *ctl.User) {
	g.GET("/users", context.WrapContext(c.Query, "查询用户数据"))
	g.GET("/users/:id", context.WrapContext(c.Get, "查询指定用户数据"))
	g.POST("/users", context.WrapContext(c.Create, "创建用户数据"))
	g.PUT("/users/:id", context.WrapContext(c.Update, "更新用户数据"))
	g.DELETE("/users/:id", context.WrapContext(c.Delete, "删除用户数据"))
	g.DELETE("/users", context.WrapContext(c.DeleteMany, "删除多条用户数据"))
	g.PATCH("/users/:id/enable", context.WrapContext(c.Enable, "启用用户数据"))
	g.PATCH("/users/:id/disable", context.WrapContext(c.Disable, "禁用用户数据"))
}
