package router

import (
	"gin-admin/src/api"
	"gin-admin/src/context"

	"github.com/gin-gonic/gin"
)

// APIUserRouter 注册/users路由
func APIUserRouter(g *gin.RouterGroup, user *api.User) {
	g.GET("/users", context.WrapContext(user.Query, "查询用户数据"))
	g.GET("/users/:id", context.WrapContext(user.Get, "查询指定用户数据"))
	g.POST("/users", context.WrapContext(user.Create, "创建用户数据"))
	g.PUT("/users/:id", context.WrapContext(user.Update, "更新用户数据"))
	g.DELETE("/users/:id", context.WrapContext(user.Delete, "删除用户数据"))
	g.DELETE("/users", context.WrapContext(user.DeleteMany, "删除多条用户数据"))
	g.PATCH("/users/:id/enable", context.WrapContext(user.Enable, "启用用户数据"))
	g.PATCH("/users/:id/disable", context.WrapContext(user.Disable, "禁用用户数据"))
}
