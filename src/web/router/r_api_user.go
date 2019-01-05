package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIUserRouter 注册/users路由
func APIUserRouter(g *gin.RouterGroup, c *ctl.User) {
	GET(g, "/users", c.Query, "查询用户数据")
	GET(g, "/users/:id", c.Get, "查询指定用户数据")
	POST(g, "/users", c.Create, "创建用户数据")
	PUT(g, "/users/:id", c.Update, "更新用户数据")
	DELETE(g, "/users/:id", c.Delete, "删除用户数据")
	DELETE(g, "/users", c.DeleteMany, "删除多条用户数据")
	PATCH(g, "/users/:id/enable", c.Enable, "启用用户数据")
	PATCH(g, "/users/:id/disable", c.Disable, "禁用用户数据")
}
