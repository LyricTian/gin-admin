package middleware

import (
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.Enforcer, skipper ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}

		ctx := context.New(c)
		p := c.Request.URL.Path
		m := c.Request.Method
		if b, err := enforcer.EnforceSafe(ctx.GetUserID(), p, m); err != nil {
			logger.StartSpan(ctx.GetContext(), "casbin中间件", "CasbinMiddleware").
				Errorf(err.Error())
			ctx.ResError(errors.NewInternalServerError())
			return
		} else if !b {
			ctx.ResError(errors.NewForbiddenError("没有访问权限"))
			return
		}
		c.Next()
	}
}
