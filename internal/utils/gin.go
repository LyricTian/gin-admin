package utils

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/LyricTian/gin-admin/v10/pkg/util/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

// Get access token from header or query parameter
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = auth
	}

	if token == "" {
		token = c.Query("accessToken")
	}

	return token
}

// Get body data from context
func GetBodyData(c *gin.Context) []byte {
	if v, ok := c.Get(RequestBodyKey); ok {
		if b, ok := v.([]byte); ok {
			return b
		}
	}
	return nil
}

// Parse body json data to struct
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.BadRequest("", "Failed to parse json: %s", err.Error())
	}
	return nil
}

// Parse query parameter to struct
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.BadRequest("", "Failed to parse query: %s", err.Error())
	}
	return nil
}

// Parse body form data to struct
func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return errors.BadRequest("", "Failed to parse form: %s", err.Error())
	}
	return nil
}

// Response json data with status code
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.Set(ResponseBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

// Response data object
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
		Data:    v,
	})
}

// Response success
func ResOK(c *gin.Context) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
	})
}

// Response pagination data
func ResPage(c *gin.Context, v interface{}, pr *PaginationResult) {
	var total int64
	if pr != nil {
		total = pr.Total
	}

	reflectValue := reflect.Indirect(reflect.ValueOf(v))
	if reflectValue.IsNil() {
		v = make([]interface{}, 0)
	}

	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
		Data:    v,
		Total:   total,
	})
}

// Response error and parse status code
func ResError(c *gin.Context, err error, status ...int) {
	ctx := c.Request.Context()

	var ierr *errors.Error
	if e, ok := errors.As(err); ok {
		ierr = e
	} else {
		ierr = errors.FromError(errors.InternalServerError("", err.Error()))
	}

	code := int(ierr.Code)
	if len(status) > 0 {
		code = status[0]
	}

	fields := []zap.Field{
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("referer", c.Request.Referer()),
		zap.String("uri", c.Request.RequestURI),
		zap.String("host", c.Request.Host),
		zap.String("remote_addr", c.Request.RemoteAddr),
		zap.String("proto", c.Request.Proto),
		zap.Int64("content_length", c.Request.ContentLength),
		zap.String("pragma", c.Request.Header.Get("Pragma")),
		zap.Int("code", code),
		zap.Error(err),
	}

	ctx = logging.NewTag(ctx, logging.TagKeySystem)
	if code >= 400 && code < 500 {
		logging.Context(ctx).Info(ierr.Detail, fields...)
	} else if code >= 500 {
		logging.Context(ctx).Error(ierr.Detail, fields...)
	}

	if code >= 500 {
		ierr.Detail = http.StatusText(http.StatusInternalServerError)
	}

	ierr.Code = int32(code)
	ResJSON(c, code, ResponseResult{Error: ierr})
}
