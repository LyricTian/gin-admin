package router

import (
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIMenuRouter 注册/menus路由
func APIMenuRouter(g *gin.RouterGroup, c *ctl.Menu) {
	g.GET("/menus", context.WrapContext(c.Query, "查询菜单数据"))
	g.GET("/menus/:id", context.WrapContext(c.Get, "查询指定菜单数据"))
	g.POST("/menus", context.WrapContext(c.Create, "创建菜单数据"))
	g.PUT("/menus/:id", context.WrapContext(c.Update, "更新菜单数据"))
	g.DELETE("/menus/:id", context.WrapContext(c.Delete, "删除菜单数据"))
	g.DELETE("/menus", context.WrapContext(c.DeleteMany, "删除多条菜单数据"))
	g.PATCH("/menus/:id/enable", context.WrapContext(c.Enable, "启用菜单数据"))
	g.PATCH("/menus/:id/disable", context.WrapContext(c.Disable, "禁用菜单数据"))
}
