package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APILoginRouter 注册登录相关路由
func APILoginRouter(g *gin.RouterGroup, c *ctl.Login) {
	GET(g, "/v1/login/captchaid", c.GetCaptchaID, SetTitle("获取验证码ID"))
	GET(g, "/v1/login/captcha", c.GetCaptcha, SetTitle("获取图形验证码"))
	POST(g, "/v1/login", c.Login, SetTitle("用户登录"))
	POST(g, "/v1/login/exit", c.Logout, SetTitle("用户登出"))
	PUT(g, "/v1/current/password", c.UpdatePassword, SetTitle("更新个人密码"))
	GET(g, "/v1/current/user", c.GetUserInfo, SetTitle("获取当前用户信息"))
	GET(g, "/v1/current/menutree", c.QueryUserMenuTree, SetTitle("查询当前用户菜单树"))
}
