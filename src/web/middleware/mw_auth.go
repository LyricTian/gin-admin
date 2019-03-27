package middleware

import (
	"github.com/LyricTian/gin-admin/src/auth"
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
)

// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(a *auth.Auth, skipper ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.New(c)

		var userID string
		if t := ctx.GetToken(); t != "" {
			id, err := a.ParseUserID(t)
			if err != nil {
				if err == auth.ErrInvalidToken {
					ctx.ResError(errors.NewUnauthorizedError())
					return
				}
				logger.StartSpan(ctx.GetContext(), "用户授权中间件", "UserAuthMiddleware").Errorf(err.Error())
				ctx.ResError(errors.NewInternalServerError())
				return
			}
			userID = id
		}

		if userID != "" {
			c.Set(context.UserIDKey, userID)
		}

		if len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}

		if userID == "" {
			if config.IsDebugMode() {
				c.Set(context.UserIDKey, config.GetRoot().UserName)
				c.Next()
				return
			}
			ctx.ResError(errors.NewUnauthorizedError("用户未登录"))
		}
	}
}
