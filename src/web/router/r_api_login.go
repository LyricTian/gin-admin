package router

import (
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APILoginRouter 注册登录相关路由
func APILoginRouter(g *gin.RouterGroup, c *ctl.Login) {
	g.POST("/login", context.WrapContext(c.Login, "用户登录"))
	g.POST("/logout", context.WrapContext(c.Logout, "用户登出"))
	g.GET("/current/user", context.WrapContext(c.GetCurrentUserInfo, "获取当前用户信息"))
	g.GET("/current/menus", context.WrapContext(c.QueryCurrentUserMenus, "查询当前用户菜单"))
}
