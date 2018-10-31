package router

import (
	"gin-admin/src/api"
	"gin-admin/src/context"

	"github.com/gin-gonic/gin"
)

// APILoginRouter 注册登录相关路由
func APILoginRouter(g *gin.RouterGroup, login *api.Login) {
	g.POST("/login", context.WrapContext(login.Login, "用户登录"))
}
