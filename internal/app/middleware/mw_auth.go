package middleware

import (
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/ginplus"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/gin-gonic/gin"
)

// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if t := ginplus.GetToken(c); t != "" {
			id, err := a.ParseUserID(ginplus.NewContext(c), t)
			if err != nil {
				if err == auth.ErrInvalidToken {
					ginplus.ResError(c, errors.ErrInvalidToken)
					return
				}

				e := errors.UnWrapResponse(errors.ErrInvalidToken)
				ginplus.ResError(c, errors.WrapResponse(err, e.Code, e.Message, e.StatusCode))
				return
			} else if id != "" {
				c.Set(ginplus.UserIDKey, id)
				c.Next()
				return
			}
		}

		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		cfg := config.Global()
		if cfg.IsDebugMode() {
			c.Set(ginplus.UserIDKey, cfg.Root.UserName)
			c.Next()
			return
		}
		ginplus.ResError(c, errors.ErrInvalidToken)
	}
}
