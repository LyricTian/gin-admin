package errors

import (
	"github.com/pkg/errors"
)

// 定义错误函数的别名
var (
	New       = errors.New
	Wrap      = errors.Wrap
	Wrapf     = errors.Wrapf
	WithStack = errors.WithStack
)

// 定义错误
var (
	// 公共错误
	ErrNotFound                = New("资源不存在")
	ErrMethodNotAllow          = New("方法不被允许")
	ErrBadRequest              = New("请求发生错误")
	ErrInvalidRequestParameter = New("无效的请求参数")
	ErrTooManyRequests         = New("请求过于频繁")
	ErrUnknownQuery            = New("未知的查询类型")
	ErrInvalidParent           = New("无效的父级节点")
	ErrNotAllowDeleteWithChild = New("含有子级，不能删除")
	ErrResourceExists          = New("资源已经存在")
	ErrResourceNotAllowDelete  = New("资源不允许删除")

	// 权限错误
	ErrNoPerm         = New("无访问权限")
	ErrNoResourcePerm = New("无资源的访问权限")

	// 用户错误
	ErrInvalidUserName = New("无效的用户名")
	ErrInvalidPassword = New("无效的密码")
	ErrInvalidUser     = New("无效的用户")
	ErrUserDisable     = New("用户被禁用")
	ErrUserNotEmptyPwd = New("密码不允许为空")

	// login
	ErrLoginNotAllowModifyPwd = New("不允许修改密码")
	ErrLoginInvalidOldPwd     = New("旧密码不正确")
	ErrLoginInvalidVerifyCode = New("无效的验证码")
)

func init() {
	// 公共错误
	newBadRequestError(ErrBadRequest)
	newBadRequestError(ErrInvalidRequestParameter)
	newErrorCode(ErrNotFound, 404, ErrNotFound.Error(), 404)
	newErrorCode(ErrMethodNotAllow, 405, ErrMethodNotAllow.Error(), 405)
	newErrorCode(ErrTooManyRequests, 429, ErrTooManyRequests.Error(), 429)
	newBadRequestError(ErrUnknownQuery)
	newBadRequestError(ErrInvalidParent)
	newBadRequestError(ErrNotAllowDeleteWithChild)
	newBadRequestError(ErrResourceExists)
	newBadRequestError(ErrResourceNotAllowDelete)

	// 权限错误
	newErrorCode(ErrNoPerm, 9999, ErrNoPerm.Error(), 401)
	newErrorCode(ErrNoResourcePerm, 401, ErrNoResourcePerm.Error(), 401)

	// 用户错误
	newBadRequestError(ErrInvalidUserName)
	newBadRequestError(ErrInvalidPassword)
	newBadRequestError(ErrInvalidUser)
	newBadRequestError(ErrUserDisable)
	newBadRequestError(ErrUserNotEmptyPwd)

	// login
	newBadRequestError(ErrLoginNotAllowModifyPwd)
	newBadRequestError(ErrLoginInvalidOldPwd)
	newBadRequestError(ErrLoginInvalidVerifyCode)
}
