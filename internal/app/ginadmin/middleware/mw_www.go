package middleware

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// WWWMiddleware 静态站点中间件
func WWWMiddleware(root string, skipper ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		fpath := filepath.Join(root, filepath.FromSlash(p))
		_, verr := os.Stat(fpath)
		if verr != nil && os.IsNotExist(verr) {
			fpath = filepath.Join(root, "index.html")
		}

		http.ServeFile(c.Writer, c.Request, fpath)
		c.Abort()
	}
}
