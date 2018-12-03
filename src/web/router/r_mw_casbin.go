package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.Enforcer, allowPrefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		allow := false
		for _, prefix := range allowPrefixes {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				allow = true
				break
			}
		}

		if !allow {
			c.Next()
			return
		}

		ctx := context.NewContext(c)
		if b, err := enforcer.EnforceSafe(ctx.GetUserID(), c.Request.URL.Path, c.Request.Method); err != nil {
			ctx.ResError(errors.Wrap(err, "验证权限发生错误"), http.StatusInternalServerError)
			return
		} else if !b {
			ctx.ResError(fmt.Errorf("没有操作权限"), http.StatusUnauthorized, 9998)
			return
		}
		c.Next()
	}
}
