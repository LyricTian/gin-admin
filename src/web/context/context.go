package context

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	gcontext "github.com/LyricTian/gin-admin/src/context"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/gin-gonic/gin"
)

// 定义上下文中的键
const (
	prefix = "github.com/LyricTian/gin-admin"
	// UserIDKey 存储上下文中的键(用户ID)
	UserIDKey = prefix + "/user_id"
	// TraceIDKey 存储上下文中的键(跟踪ID)
	TraceIDKey = prefix + "/trace_id"
	// ResBodyKey 存储上下文中的键(响应Body数据)
	ResBodyKey = prefix + "/res_body"
)

// JoinRouter 拼接路由
func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}

// New 创建上下文实例
func New(c *gin.Context) *Context {
	return &Context{c}
}

// Context 定义上下文
type Context struct {
	gctx *gin.Context
}

// Reset context
func (a *Context) Reset(c *gin.Context) {
	a.gctx = c
}

func (a *Context) getFuncName(name string) string {
	return fmt.Sprintf("web.context.Context.%s", name)
}

// GetToken 获取用户令牌
func (a *Context) GetToken() string {
	var token string
	auth := a.gctx.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// GetContext get context.Context
func (a *Context) GetContext() context.Context {
	parent := context.Background()

	if v := a.GetTraceID(); v != "" {
		parent = logger.NewTraceIDContext(parent, a.GetTraceID())
	}

	if v := a.GetUserID(); v != "" {
		parent = gcontext.NewUserID(parent, v)
		parent = logger.NewUserIDContext(parent, v)
	}

	return parent
}

// Request http request
func (a *Context) Request() *http.Request {
	return a.gctx.Request
}

// ResponseWriter http response stream
func (a *Context) ResponseWriter() http.ResponseWriter {
	return a.gctx.Writer
}

// Param 获取路径参数(/foo/:id)
func (a *Context) Param(key string) string {
	return a.gctx.Param(key)
}

// Query 获取查询参数(/foo?id=)
func (a *Context) Query(key string) string {
	return a.gctx.Query(key)
}

// GetPageIndex 获取分页的页索引
func (a *Context) GetPageIndex() int {
	if v := a.Query("current"); v != "" {
		if iv := util.S(v).Int(); iv > 0 {
			return iv
		}
	}
	return 1
}

// GetPageSize 获取分页的页大小(最大50)
func (a *Context) GetPageSize() int {
	if v := a.Query("pageSize"); v != "" {
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
func (a *Context) GetPaginationParam() *schema.PaginationParam {
	return &schema.PaginationParam{
		PageIndex: a.GetPageIndex(),
		PageSize:  a.GetPageSize(),
	}
}

// GetTraceID 获取追踪ID
func (a *Context) GetTraceID() string {
	return a.gctx.GetString(TraceIDKey)
}

// GetUserID 获取用户ID
func (a *Context) GetUserID() string {
	return a.gctx.GetString(UserIDKey)
}

// SetUserID 设定用户ID
func (a *Context) SetUserID(userID string) {
	a.gctx.Set(UserIDKey, userID)
}

// ParseJSON 解析请求JSON
func (a *Context) ParseJSON(obj interface{}) error {
	if err := a.gctx.ShouldBindJSON(obj); err != nil {
		logger.StartSpan(a.GetContext(), "解析请求JSON", a.getFuncName("ParseJSON")).Warnf(err.Error())
		return errors.NewBadRequestError("无效的请求参数")
	}
	return nil
}

// 根据错误获取状态码
func (a *Context) getStatusByError(err error, status int) int {
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
func (a *Context) ResError(err error, code ...int) {
	a.ResErrorWithStatus(err, 0, code...)
}

// ResErrorWithStatus 响应错误和指定状态码(不指定则根据错误自动判断)
func (a *Context) ResErrorWithStatus(err error, status int, code ...int) {
	var item schema.HTTPErrorItem

	switch e := err.(type) {
	case *errors.MessageError:
		item.Message = e.Error()
		status = a.getStatusByError(e.Parent(), status)
	default:
		if err != nil {
			item.Message = err.Error()
		}
		status = a.getStatusByError(err, status)
	}

	if len(code) > 0 {
		item.Code = code[0]
	}

	a.ResJSON(status, schema.HTTPError{Error: item})
}

// ResPage 响应分页数据
func (a *Context) ResPage(v interface{}, pr *schema.PaginationResult) {
	list := schema.HTTPList{
		List: v,
		Pagination: &schema.HTTPPagination{
			Current:  a.GetPageIndex(),
			PageSize: a.GetPageSize(),
		},
	}
	if pr != nil {
		list.Pagination.Total = pr.Total
	}

	a.ResSuccess(list)
}

// ResList 响应列表数据
func (a *Context) ResList(v interface{}) {
	a.ResSuccess(schema.HTTPList{List: v})
}

// ResOK 响应OK
func (a *Context) ResOK() {
	a.ResSuccess(schema.HTTPStatus{Status: "OK"})
}

// ResSuccess 响应成功
func (a *Context) ResSuccess(v interface{}) {
	a.ResJSON(http.StatusOK, v)
}

// ResJSON 响应JSON数据
func (a *Context) ResJSON(status int, v interface{}) {
	buf, err := util.JSONMarshal(v)
	if err != nil {
		logger.StartSpan(a.GetContext(), "响应JSON数据", a.getFuncName("ResJSON")).
			WithField("object", v).Errorf(err.Error())
		a.ResError(errors.NewInternalServerError())
		return
	}
	a.gctx.Set(ResBodyKey, buf)
	a.gctx.Data(status, "application/json; charset=utf-8", buf)
}
