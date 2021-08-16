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
	ErrInvalidToken            = NewResponse(40001, 401, "令牌失效")
	ErrNoPerm                  = NewResponse(0, 401, "无访问权限")
	ErrNotFound                = NewResponse(0, 404, "资源不存在")
	ErrMethodNotAllow          = NewResponse(0, 405, "方法不被允许")
	ErrTooManyRequests         = NewResponse(0, 429, "请求过于频繁")
	ErrInternalServer          = NewResponse(0, 500, "服务器发生错误")
	ErrBadRequest              = New400Response("请求发生错误")
	ErrInvalidUser             = New400Response("无效的用户")
	ErrUserDisable             = New400Response("用户被禁用，请联系管理员")
	ErrInvalidParent           = New400Response("Not found parent node")
	ErrNotAllowDeleteWithChild = New400Response("存在子节点，不允许删除")
)
