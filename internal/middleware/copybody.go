package middleware

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"

	"github.com/gin-gonic/gin"
)

type CopyBodyConfig struct {
	SkippedPathPrefixes []string
	MaxContentLen       int64
}

var DefaultCopyBodyConfig = CopyBodyConfig{
	MaxContentLen: 32 << 20, // 32MB
}

func CopyBody() gin.HandlerFunc {
	return CopyBodyWithConfig(DefaultCopyBodyConfig)
}

func CopyBodyWithConfig(config CopyBodyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkippedPathPrefixes(c, config.SkippedPathPrefixes...) || c.Request.Body == nil {
			c.Next()
			return
		}

		var (
			requestBody []byte
			err         error
		)

		isGzip := false
		safe := http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxContentLen)
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
			utilx.ResError(c, errors.RequestEntityTooLarge(errors.ErrRequestEntityTooLarge,
				"Request body too large, limit %d byte", config.MaxContentLen))
			return
		}

		c.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		c.Request.Body = ioutil.NopCloser(bf)
		c.Set(utilx.RequestBodyKey, requestBody)

		c.Next()
	}
}
