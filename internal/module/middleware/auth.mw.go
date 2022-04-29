package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/module/consts"
	"github.com/LyricTian/gin-admin/v9/internal/module/contextx"
	"github.com/LyricTian/gin-admin/v9/internal/module/ginx"
	"github.com/LyricTian/gin-admin/v9/internal/module/util"
	"github.com/LyricTian/gin-admin/v9/pkg/cache"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/jwtauth"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
)

func wrapUserAuthContext(c *gin.Context, subject string) {
	sub := util.ParseSubject(subject)
	ctx := c.Request.Context()

	ctx = contextx.NewUserID(ctx, sub.UserID)
	ctx = contextx.NewRole(ctx, sub.Role)
	ctx = logger.NewUserIDContext(ctx, sub.UserID)

	c.Request = c.Request.WithContext(ctx)
}

// Verification JWT middleware
func UserAuthMiddleware(a jwtauth.Auther, cache cache.Cacher, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return func(c *gin.Context) {
			wrapUserAuthContext(c, config.C.Root.UserID)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		token := ginx.GetToken(c)
		if token == "" && SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		subject, err := a.ParseSubject(c.Request.Context(), token)
		if err != nil {
			if err == jwtauth.ErrInvalidToken {
				if config.C.IsDebugMode() {
					wrapUserAuthContext(c, util.GenerateSubject(util.Subject{
						UserID: config.C.Root.UserID,
						Role:   "Root",
					}))
					c.Next()
					return
				}
				ginx.ResError(c, errors.Unauthorized(consts.ErrInvalidTokenID, err.Error()))
				return
			}
			ginx.ResError(c, errors.InternalServerError(consts.ErrInternalServerErrorID, err.Error()))
			return
		}

		wrapUserAuthContext(c, subject)

		// Validate user_id exists in cache or database
		ctx := c.Request.Context()
		userID := contextx.FromUserID(ctx)

		// If user_id is root, skip cache
		if config.C.Root.UserID == userID {
			c.Next()
			return
		}

		c.Next()
	}
}
