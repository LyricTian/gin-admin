package middleware

import (
	"fmt"

	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.Enforcer, skipPrefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.New(c)

		p := c.Request.URL.Path
		m := c.Request.Method

		// 跳过不需要校验权限的路由(规则：GET/api/v1/test)
		if util.CheckPrefix(fmt.Sprintf("%s%s", p, m), skipPrefixes...) {
			c.Next()
			return
		}

		if b, err := enforcer.EnforceSafe(ctx.GetUserID(), p, m); err != nil {
			logger.Start(ctx.CContext()).Errorf("验证权限发生错误: %s", err.Error())
			ctx.ResError(util.NewInternalServerError("服务器发生权限校验错误"))
			return
		} else if !b {
			ctx.ResError(util.NewUnauthorizedError("无操作权限"), 9998)
			return
		}
		c.Next()
	}
}
