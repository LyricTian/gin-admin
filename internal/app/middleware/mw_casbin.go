package middleware

import (
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/ginplus"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.Enforcer, skipper ...SkipperFunc) gin.HandlerFunc {
	cfg := config.GetGlobalConfig()
	return func(c *gin.Context) {
		if !cfg.EnableCasbin || len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		m := c.Request.Method
		if b, err := enforcer.EnforceSafe(ginplus.GetUserID(c), p, m); err != nil {
			ginplus.ResError(c, errors.WithStack(err))
			return
		} else if !b {
			ginplus.ResError(c, errors.ErrNoResourcePerm)
			return
		}
		c.Next()
	}
}
