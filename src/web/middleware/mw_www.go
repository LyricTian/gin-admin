package middleware

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/LyricTian/gin-admin/src/util"
	"github.com/gin-gonic/gin"
)

// WWWMiddleware 静态站点中间件
func WWWMiddleware(root string, excludeRouters ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		m := c.Request.Method
		p := c.Request.URL.Path
		if util.CheckPrefix(p, excludeRouters...) ||
			!(m == http.MethodHead || m == http.MethodGet) {
			c.Next()
			return
		}

		fpath := filepath.Join(root, filepath.FromSlash(p))
		_, verr := os.Stat(fpath)
		if verr != nil && os.IsNotExist(verr) {
			fpath = filepath.Join(root, "index.html")
		}

		http.ServeFile(c.Writer, c.Request, fpath)
		c.Abort()
	}
}
