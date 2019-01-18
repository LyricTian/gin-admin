package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIMenuRouter 注册/menus路由
func APIMenuRouter(g *gin.RouterGroup, c *ctl.Menu) {
	GET(g, "/v1/menus", c.Query, SetTitle("查询菜单数据"))
	GET(g, "/v1/menus/:id", c.Get, SetTitle("查询指定菜单数据"))
	POST(g, "/v1/menus", c.Create, SetTitle("创建菜单数据"))
	PUT(g, "/v1/menus/:id", c.Update, SetTitle("更新菜单数据"))
	DELETE(g, "/v1/menus/:id", c.Delete, SetTitle("删除菜单数据"))
	DELETE(g, "/v1/menus", c.DeleteMany, SetTitle("删除多条菜单数据"))
}
