package utilx

import (
	"net/http"
	"strings"

	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/LyricTian/gin-admin/v9/pkg/util/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

// Get jwt token from header (Authorization: Bearer xxx)
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
		return errors.BadRequest(errors.ErrBadRequestID, "Failed to parse json: %s", err.Error())
	}
	return nil
}

// Parse query parameter to struct
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.BadRequest(errors.ErrBadRequestID, "Failed to parse query: %s", err.Error())
	}
	return nil
}

// Parse body form data to struct
func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return errors.BadRequest(errors.ErrBadRequestID, "Failed to parse form: %s", err.Error())
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
	ResJSON(c, http.StatusOK, ResponseResult{
		Success:    true,
		Data:       v,
		Pagination: pr,
	})
}

// Response error and parse status code
func ResError(c *gin.Context, err error, status ...int) {
	ctx := c.Request.Context()

	var ierr *errors.Error
	if e, ok := errors.As(err); ok {
		ierr = e
	} else {
		ierr = errors.FromError(errors.InternalServerError(errors.ErrInternalServerErrorID, err.Error()))
	}

	code := int(ierr.Code)
	if len(status) > 0 {
		code = status[0]
	}

	if code >= 400 && code < 500 {
		logger.Context(ctx).Info(ierr.Detail, zap.Int("code", code), zap.Error(err))
	} else if code >= 500 {
		logger.Context(ctx).Error(ierr.Detail, zap.Int("code", code), zap.Error(err))
	}

	if code >= 500 {
		ierr.Detail = http.StatusText(http.StatusInternalServerError)
	}

	ResJSON(c, code, ResponseResult{Error: ierr})
}
