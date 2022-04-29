package middleware

import (
	"net/http"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/module/consts"
	"github.com/LyricTian/gin-admin/v9/internal/module/contextx"
	"github.com/LyricTian/gin-admin/v9/internal/module/ginx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// Valid use interface permission
func CasbinMiddleware(enforcer *casbin.Enforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config.C.Casbin
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		m := c.Request.Method
		if config.C.CORS.Enable && m == http.MethodOptions {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		role := contextx.FromRole(c.Request.Context())
		if b, err := enforcer.Enforce(role, p, m); err != nil {
			ginx.ResError(c, errors.InternalServerError(consts.ErrInternalServerErrorID, err.Error()))
			return
		} else if !b {
			ginx.ResError(c, errors.Unauthorized("com.perm.unauthorized", "role %s not allowed to access this api", role))
			return
		}
		c.Next()
	}
}
