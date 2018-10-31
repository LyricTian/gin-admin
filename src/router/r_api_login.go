package router

import (
	"gin-admin/src/api"
	"gin-admin/src/context"

	"github.com/gin-gonic/gin"
)

// APILoginRouter 注册登录相关路由
func APILoginRouter(g *gin.RouterGroup, login *api.Login) {
	g.POST("/login", context.WrapContext(login.Login, "用户登录"))
	g.POST("/logout", context.WrapContext(login.Logout, "用户登出"))
	g.GET("/current/user", context.WrapContext(login.GetCurrentUserInfo, "获取当前用户信息"))
}
