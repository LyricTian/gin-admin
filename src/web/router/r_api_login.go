package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APILoginRouter 注册登录相关路由
func APILoginRouter(g *gin.RouterGroup, c *ctl.Login) {
	POST(g, "/v1/login", c.Login, SetTitle("用户登录"))
	POST(g, "/v1/logout", c.Logout, SetTitle("用户登出"))
	GET(g, "/v1/current/user", c.GetCurrentUserInfo, SetTitle("获取当前用户信息"))
	GET(g, "/v1/current/menus", c.QueryCurrentUserMenus, SetTitle("查询当前用户菜单"))
}
