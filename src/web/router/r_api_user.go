package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIUserRouter 注册/users路由
func APIUserRouter(g *gin.RouterGroup, c *ctl.User) {
	GET(g, "/v1/users", c.Query, SetName("查询用户数据"))
	GET(g, "/v1/users/:id", c.Get, SetName("查询指定用户数据"))
	POST(g, "/v1/users", c.Create, SetName("创建用户数据"))
	PUT(g, "/v1/users/:id", c.Update, SetName("更新用户数据"))
	DELETE(g, "/v1/users/:id", c.Delete, SetName("删除用户数据"))
	DELETE(g, "/v1/users", c.DeleteMany, SetName("删除多条用户数据"))
	PATCH(g, "/v1/users/:id/enable", c.Enable, SetName("启用用户数据"))
	PATCH(g, "/v1/users/:id/disable", c.Disable, SetName("禁用用户数据"))
}
