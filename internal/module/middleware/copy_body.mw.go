package middleware

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/module/consts"
	"github.com/LyricTian/gin-admin/v9/internal/module/ginx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Copy body to context bytes array
func CopyBodyMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	var maxMemory int64 = 16 << 20 // 16 MB
	if v := config.C.HTTP.MaxContentLength; v > 0 {
		maxMemory = v
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) || c.Request.Body == nil {
			c.Next()
			return
		}

		var (
			requestBody []byte
			err         error
		)

		isGzip := false
		safe := http.MaxBytesReader(c.Writer, c.Request.Body, maxMemory)

		if c.GetHeader("Content-Encoding") == "gzip" {
			if reader, ierr := gzip.NewReader(safe); ierr == nil {
				isGzip = true
				requestBody, err = ioutil.ReadAll(reader)
			}
		}

		if !isGzip {
			requestBody, err = ioutil.ReadAll(safe)
		}

		if err != nil {
			ginx.ResError(c, errors.RequestEntityTooLarge(consts.ErrRequestEntityTooLarge,
				"request body too large, limit %d byte", maxMemory))
			return
		}

		c.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		c.Request.Body = ioutil.NopCloser(bf)
		c.Set(ginx.ReqBodyKey, requestBody)

		c.Next()
	}
}
