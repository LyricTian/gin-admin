package middleware

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// WWWMiddleware 静态站点中间件
func WWWMiddleware(root string, skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		fpath := filepath.Join(root, filepath.FromSlash(p))
		_, err := os.Stat(fpath)
		if err != nil && os.IsNotExist(err) {
			fpath = filepath.Join(root, "index.html")
		}

		http.ServeFile(c.Writer, c.Request, fpath)
		c.Abort()
	}
}
