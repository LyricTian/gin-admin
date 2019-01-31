package middleware

import (
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/web/auth"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
)

// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(skipper SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipper != nil && skipper(c) {
			c.Next()
			return
		}

		ctx := context.New(c)
		userInfo, err := auth.GetUserInfo(c)
		if err != nil {
			ctx.ResError(err)
			return
		} else if userInfo == nil {
			if !config.IsDebugMode() {
				ctx.ResError(errors.NewUnauthorizedError("用户未登录"))
				return
			}
			// 调试模式下，如果用户未登录，则使用root用户
			userInfo = &auth.UserInfo{
				UserID: config.GetRootUser().UserName,
			}
		}

		c.Set(context.UserIDKey, userInfo.UserID)
		c.Next()
	}
}
