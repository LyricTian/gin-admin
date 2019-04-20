package ginplus

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	icontext "github.com/LyricTian/gin-admin/internal/app/ginadmin/context"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/gin-gonic/gin"
)

// 定义上下文中的键
const (
	prefix = "ginadmin"
	// UserIDKey 存储上下文中的键(用户ID)
	UserIDKey = prefix + "/user_id"
	// TraceIDKey 存储上下文中的键(跟踪ID)
	TraceIDKey = prefix + "/trace_id"
	// ResBodyKey 存储上下文中的键(响应Body数据)
	ResBodyKey = prefix + "/res_body"
)

func getFuncName(name string) string {
	return fmt.Sprintf("ginadmin.ginplus.%s", name)
}

// NewContext get context.Context
func NewContext(c *gin.Context) context.Context {
	parent := context.Background()

	if v := GetTraceID(c); v != "" {
		parent = icontext.NewTraceID(parent, v)
		parent = logger.NewTraceIDContext(parent, GetTraceID(c))
	}

	if v := GetUserID(c); v != "" {
		parent = icontext.NewUserID(parent, v)
		parent = logger.NewUserIDContext(parent, v)
	}

	return parent
}

// GetToken 获取用户令牌
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// GetPageIndex 获取分页的页索引
func GetPageIndex(c *gin.Context) int {
	if v := c.Query("current"); v != "" {
		if iv := util.S(v).Int(); iv > 0 {
			return iv
		}
	}
	return 1
}

// GetPageSize 获取分页的页大小(最大50)
func GetPageSize(c *gin.Context) int {
	if v := c.Query("pageSize"); v != "" {
		if iv := util.S(v).Int(); iv > 0 {
			if iv > 50 {
				iv = 50
			}
			return iv
		}
	}
	return 10
}

// GetPaginationParam 获取分页查询参数
func GetPaginationParam(c *gin.Context) *schema.PaginationParam {
	return &schema.PaginationParam{
		PageIndex: GetPageIndex(c),
		PageSize:  GetPageSize(c),
	}
}

// GetTraceID 获取追踪ID
func GetTraceID(c *gin.Context) string {
	return c.GetString(TraceIDKey)
}

// GetUserID 获取用户ID
func GetUserID(c *gin.Context) string {
	return c.GetString(UserIDKey)
}

// SetUserID 设定用户ID
func SetUserID(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)
}

// ParseJSON 解析请求JSON
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		logger.StartSpan(NewContext(c), "解析请求JSON", getFuncName("ParseJSON")).Warnf(err.Error())
		return errors.NewBadRequestError("无效的请求参数")
	}
	return nil
}

// 根据错误获取状态码
func getStatusByError(c *gin.Context, err error, status int) int {
	if status > 0 {
		return status
	}

	switch err {
	case errors.ErrBadRequest:
		status = 400
	case errors.ErrUnauthorized:
		status = 401
	case errors.ErrForbidden:
		status = 403
	case errors.ErrNotFound:
		status = 404
	case errors.ErrInternalServer:
		status = 500
	default:
		status = 500
	}
	return status
}

// ResError 响应错误
func ResError(c *gin.Context, err error, code ...int) {
	ResErrorWithStatus(c, err, 0, code...)
}

// ResErrorWithStatus 响应错误和指定状态码(不指定则根据错误自动判断)
func ResErrorWithStatus(c *gin.Context, err error, status int, code ...int) {
	var item schema.HTTPErrorItem

	switch e := err.(type) {
	case *errors.MessageError:
		item.Message = e.Error()
		status = getStatusByError(c, e.Parent(), status)
	default:
		if err != nil {
			item.Message = err.Error()
		}
		status = getStatusByError(c, err, status)
	}

	if len(code) > 0 {
		item.Code = code[0]
	}

	ResJSON(c, status, schema.HTTPError{Error: item})
}

// ResPage 响应分页数据
func ResPage(c *gin.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.HTTPList{
		List: v,
		Pagination: &schema.HTTPPagination{
			Current:  GetPageIndex(c),
			PageSize: GetPageSize(c),
		},
	}
	if pr != nil {
		list.Pagination.Total = pr.Total
	}

	ResSuccess(c, list)
}

// ResList 响应列表数据
func ResList(c *gin.Context, v interface{}) {
	ResSuccess(c, schema.HTTPList{List: v})
}

// ResOK 响应OK
func ResOK(c *gin.Context) {
	ResSuccess(c, schema.HTTPStatus{Status: "OK"})
}

// ResSuccess 响应成功
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

// ResJSON 响应JSON数据
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := util.JSONMarshal(v)
	if err != nil {
		logger.StartSpan(NewContext(c), "响应JSON数据", getFuncName("ResJSON")).
			WithField("object", v).Errorf(err.Error())
		ResError(c, errors.NewInternalServerError())
		return
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}
