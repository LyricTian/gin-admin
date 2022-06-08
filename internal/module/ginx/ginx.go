package ginx

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"

	"github.com/LyricTian/gin-admin/v9/internal/module/consts"
	"github.com/LyricTian/gin-admin/v9/internal/schema"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/LyricTian/gin-admin/v9/pkg/util/json"
)

const (
	prefix     = "ginadmin"
	ReqBodyKey = prefix + "/req-body"
	ResBodyKey = prefix + "/res-body"
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
	if v, ok := c.Get(ReqBodyKey); ok {
		if b, ok := v.([]byte); ok {
			return b
		}
	}
	return nil
}

// Parse body json data to struct
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.BadRequest(consts.ErrBadRequestID, "Parse json failed: %s", err.Error())
	}
	return nil
}

// Parse query parameter to struct
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.BadRequest(consts.ErrBadRequestID, "Parse query failed: %s", err.Error())
	}
	return nil
}

// Parse body form data to struct
func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return errors.BadRequest(consts.ErrBadRequestID, "Parse form failed: %s", err.Error())
	}
	return nil
}

// Response success with status ok
func ResOK(c *gin.Context) {
	ResSuccess(c, schema.OkResult{Ok: true})
}

// Response data with list object
func ResList(c *gin.Context, v interface{}) {
	ResSuccess(c, schema.ListResult{List: v})
}

// Response pagination data object
func ResPage(c *gin.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.ListResult{
		List:       v,
		Pagination: pr,
	}
	ResSuccess(c, list)
}

// Response data object
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

// Response json data with status code
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

// Response error object and parse error status code
func ResError(c *gin.Context, err error, status ...int) {
	ctx := c.Request.Context()

	var ierr *errors.Error
	if e, ok := errors.As(err); ok {
		ierr = e
	} else {
		ierr = errors.FromError(errors.InternalServerError(consts.ErrInternalServerErrorID, err.Error()))
	}

	code := int(ierr.Code)
	if len(status) > 0 {
		code = status[0]
	}

	if code >= 400 && code < 500 {
		logger.WithContext(ctx).WithFields(logrus.Fields{
			"error_id":     ierr.ID,
			"error_code":   ierr.Code,
			"error_status": ierr.Status,
		}).Info(ierr.Detail)
	} else if code >= 500 {
		logger.WithContext(logger.NewStackContext(ctx, err)).
			WithFields(logrus.Fields{
				"error_id":     ierr.ID,
				"error_code":   ierr.Code,
				"error_status": ierr.Status,
			}).Error(ierr.Detail)
	}

	if code >= 500 {
		ierr.Detail = "Internal server error"
	}

	ResJSON(c, code, schema.ErrorResult{Error: ierr})
}
