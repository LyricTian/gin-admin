package middleware

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type StaticConfig struct {
	SkippedPathPrefixes []string
	Root                string
}

func StaticWithConfig(config StaticConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkippedPathPrefixes(c, config.SkippedPathPrefixes...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		fpath := filepath.Join(config.Root, filepath.FromSlash(p))
		_, err := os.Stat(fpath)
		if err != nil && os.IsNotExist(err) {
			fpath = filepath.Join(config.Root, "index.html")
		}
		c.File(fpath)
		c.Abort()
	}
}
