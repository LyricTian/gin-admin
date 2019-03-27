package errors

import (
	"github.com/pkg/errors"
)

// 定义通用错误
var (
	New               = errors.New
	Wrap              = errors.Wrap
	Wrapf             = errors.Wrapf
	ErrForbidden      = New("禁止访问")
	ErrNotFound       = New("资源不存在")
	ErrBadRequest     = New("请求无效")
	ErrUnauthorized   = New("未授权")
	ErrInternalServer = New("服务器错误")
)

// NewBadRequestError 创建请求无效错误
func NewBadRequestError(msg ...string) error {
	return NewMessageError(ErrBadRequest, msg...)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(msg ...string) error {
	return NewMessageError(ErrUnauthorized, msg...)
}

// NewForbiddenError 创建资源禁止访问错误
func NewForbiddenError(msg ...string) error {
	return NewMessageError(ErrForbidden, msg...)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(msg ...string) error {
	return NewMessageError(ErrNotFound, msg...)
}

// NewInternalServerError 创建服务器错误
func NewInternalServerError(msg ...string) error {
	return NewMessageError(ErrInternalServer, msg...)
}

// NewMessageError 创建自定义消息错误
func NewMessageError(parent error, msg ...string) error {
	if parent == nil {
		return nil
	}

	m := parent.Error()
	if len(msg) > 0 {
		m = msg[0]
	}
	return &MessageError{parent, m}
}

// MessageError 自定义消息错误
type MessageError struct {
	err error
	msg string
}

func (m *MessageError) Error() string {
	return m.msg
}

// Parent 父级错误
func (m *MessageError) Parent() error {
	return m.err
}
