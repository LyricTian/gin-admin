package middlewares

import (
	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/gin-gonic/gin"
)

type AuthConfig struct {
	SkippedPathPrefixes []string
	Skipper             func(c *gin.Context) bool
	DefaultUserID       func(c *gin.Context) string
	ParseUserID         func(c *gin.Context) (string, error)
}

func AuthWithConfig(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkippedPathPrefixes(c, config.SkippedPathPrefixes...) ||
			(config.Skipper != nil && config.Skipper(c)) {
			c.Next()
			return
		}

		userID, err := config.ParseUserID(c)
		if err != nil {
			utils.ResError(c, err)
			return
		}

		ctx := utils.NewUserID(c.Request.Context(), userID)
		if userID == config.DefaultUserID(c) {
			ctx = utils.NewIsRootUser(ctx)
		}
		ctx = logging.NewUserID(ctx, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
