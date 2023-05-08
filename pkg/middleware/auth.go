package middleware

import (
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"github.com/gin-gonic/gin"
)

type AuthConfig struct {
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
	Skipper             func(c *gin.Context) bool
	ParseUserID         func(c *gin.Context) (string, error)
}

func AuthWithConfig(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			SkippedPathPrefixes(c, config.SkippedPathPrefixes...) ||
			(config.Skipper != nil && config.Skipper(c)) {
			c.Next()
			return
		}

		userID, err := config.ParseUserID(c)
		if err != nil {
			util.ResError(c, err)
			return
		}

		ctx := util.NewUserID(c.Request.Context(), userID)
		ctx = logging.NewUserID(ctx, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
