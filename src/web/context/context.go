package context

import (
	"context"
	"net/http"

	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	"github.com/go-session/session"
)

// New 创建上下文实例
func New(c *gin.Context) *Context {
	return &Context{c}
}

// Context 定义上下文
type Context struct {
	gctx *gin.Context
}

// GContext 获取gin.Context
func (a *Context) GContext() *gin.Context {
	return a.gctx
}

// CContext 获取context.Context
func (a *Context) CContext() context.Context {
	parent := context.Background()

	if v := a.GetTraceID(); v != "" {
		parent = logger.NewTraceIDContext(parent, a.GetTraceID())
	}

	if v := a.GetUserID(); v != "" {
		parent = logger.NewUserIDContext(parent, a.GetUserID())
	}

	return parent
}

// SessionStore 获取会话存储
func (a *Context) SessionStore() session.Store {
	return ginsession.FromContext(a.gctx)
}

// RefreshSession 更新会话
func (a *Context) RefreshSession() (session.Store, error) {
	return ginsession.Refresh(a.gctx)
}

// DestroySession 销毁会话
func (a *Context) DestroySession() error {
	return ginsession.Destroy(a.gctx)
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
func (a *Context) GetPageIndex() uint {
	if v := a.Query("current"); v != "" {
		if iv := util.S(v).Uint(); iv > 0 {
			return iv
		}
	}
	return 1
}

// GetPageSize 获取分页的页大小(最大50)
func (a *Context) GetPageSize() uint {
	if v := a.Query("pageSize"); v != "" {
		if iv := util.S(v).Uint(); iv > 0 {
			if iv > 50 {
				iv = 50
			}
			return iv
		}
	}
	return 10
}

// GetTraceID 获取追踪ID
func (a *Context) GetTraceID() string {
	return a.gctx.GetString(util.ContextKeyTraceID)
}

// GetUserID 获取用户ID
func (a *Context) GetUserID() string {
	return a.gctx.GetString(util.ContextKeyUserID)
}

// SetUserID 设定用户ID
func (a *Context) SetUserID(userID string) {
	a.gctx.Set(util.ContextKeyUserID, userID)
}

// ParseJSON 解析请求JSON
func (a *Context) ParseJSON(obj interface{}) error {
	if err := a.gctx.ShouldBindJSON(obj); err != nil {
		logger.Start(a.CContext()).Warnf("解析请求参数发生错误: %s", err.Error())
		return util.NewBadRequestError("解析请求参数发生错误")
	}
	return nil
}

// 根据错误获取状态码
func (a *Context) getStatusByError(err error, status int) int {
	if status > 0 {
		return status
	}

	switch err {
	case util.ErrBadRequest:
		status = 400
	case util.ErrUnauthorized:
		status = 401
	case util.ErrNotFound:
		status = 404
	case util.ErrInternalServer:
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
	var item HTTPErrorItem

	switch e := err.(type) {
	case *util.MessageError:
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
	a.ResJSON(status, HTTPError{Error: item})
}

// ResPage 响应分页数据
func (a *Context) ResPage(total int, v interface{}) {
	a.ResSuccess(HTTPList{
		List: v,
		Pagination: &HTTPPagination{
			Total:    total,
			Current:  a.GetPageIndex(),
			PageSize: a.GetPageSize(),
		},
	})
}

// ResList 响应列表数据
func (a *Context) ResList(v interface{}) {
	a.ResSuccess(HTTPList{List: v})
}

// ResOK 响应OK
func (a *Context) ResOK() {
	a.ResSuccess(HTTPStatus{Status: StatusOK})
}

// ResSuccess 响应成功
func (a *Context) ResSuccess(v interface{}) {
	if v == nil {
		v = gin.H{}
	}
	a.ResJSON(http.StatusOK, v)
}

// ResJSON 响应JSON数据
func (a *Context) ResJSON(status int, v interface{}) {
	buf, err := util.JSONMarshal(v)
	if err != nil {
		logger.Start(a.CContext()).WithField("object", v).Errorf("JSON序列化发生错误: %s", err.Error())
		a.ResError(util.NewInternalServerError())
		return
	}
	a.gctx.Set(util.ContextKeyResBody, buf)
	a.gctx.Data(status, "application/json; charset=utf-8", buf)
	a.gctx.Abort()
}
