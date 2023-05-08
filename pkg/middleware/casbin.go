package middleware

import (
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type CasbinConfig struct {
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
	Skipper             func(c *gin.Context) bool
	GetEnforcer         func(c *gin.Context) *casbin.Enforcer
	GetSubjects         func(c *gin.Context) []string
}

func CasbinWithConfig(config CasbinConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			SkippedPathPrefixes(c, config.SkippedPathPrefixes...) ||
			(config.Skipper != nil && config.Skipper(c)) {
			c.Next()
			return
		}

		enforcer := config.GetEnforcer(c)
		for _, sub := range config.GetSubjects(c) {
			if b, err := enforcer.Enforce(sub, c.Request.URL.Path, c.Request.Method); err != nil {
				util.ResError(c, err)
				return
			} else if b {
				c.Next()
				return
			}
		}
		util.ResError(c, errors.Unauthorized("com.perm.denied", "Permission denied, please contact the administrator"))
	}
}
