package router

import (
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APILoginRouter 注册登录相关路由
func APILoginRouter(g *gin.RouterGroup, login *ctl.Login) {
	g.POST("/login", context.WrapContext(login.Login, "用户登录"))
	g.POST("/logout", context.WrapContext(login.Logout, "用户登出"))
	g.GET("/current/user", context.WrapContext(login.GetCurrentUserInfo, "获取当前用户信息"))
	g.GET("/current/menus", context.WrapContext(login.QueryCurrentUserMenus, "查询当前用户菜单"))
}
