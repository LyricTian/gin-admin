package errors

import (
	"github.com/pkg/errors"
)

// 定义别名
var (
	New          = errors.New
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

// 定义错误
var (
	ErrBadRequest              = New400Response("请求发生错误")
	ErrInvalidParent           = New400Response("无效的父级节点")
	ErrNotAllowDeleteWithChild = New400Response("含有子级，不能删除")
	ErrNotAllowDelete          = New400Response("资源不允许删除")
	ErrInvalidUserName         = New400Response("无效的用户名")
	ErrInvalidPassword         = New400Response("无效的密码")
	ErrInvalidUser             = New400Response("无效的用户")
	ErrUserDisable             = New400Response("用户被禁用，请联系管理员")

	ErrNoPerm          = NewResponse(401, "无访问权限", 401)
	ErrInvalidToken    = NewResponse(9999, "令牌失效", 401)
	ErrNotFound        = NewResponse(404, "资源不存在", 404)
	ErrMethodNotAllow  = NewResponse(405, "方法不被允许", 405)
	ErrTooManyRequests = NewResponse(429, "请求过于频繁", 429)
	ErrInternalServer  = NewResponse(500, "服务器发生错误", 500)
)
