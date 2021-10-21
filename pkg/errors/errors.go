package errors

import (
	"github.com/pkg/errors"
)

// Define alias
var (
	New          = errors.New
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

var (
	ErrInvalidToken    = NewResponse(9999, 401, "invalid signature")
	ErrNoPerm          = NewResponse(0, 401, "no permission")
	ErrNotFound        = NewResponse(0, 404, "not found")
	ErrMethodNotAllow  = NewResponse(0, 405, "method not allowed")
	ErrTooManyRequests = NewResponse(0, 429, "too many requests")
	ErrInternalServer  = NewResponse(0, 500, "internal server error")
	ErrBadRequest      = New400Response("bad request")
	ErrInvalidParent   = New400Response("not found parent node")
	ErrUserDisable     = New400Response("user forbidden")
)
