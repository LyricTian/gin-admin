package errors

// ResponseError 定义响应错误
type ResponseError struct {
	Code       int    // 错误码
	Message    string // 错误消息
	StatusCode int    // 响应状态码
	ERR        error  // 响应错误
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return r.Message
}

// UnWrapResponse 解包响应错误
func UnWrapResponse(err error) *ResponseError {
	if v, ok := err.(*ResponseError); ok {
		return v
	}
	return nil
}

// WrapResponse 包装响应错误
func WrapResponse(err error, code int, msg string, status ...int) error {
	res := &ResponseError{
		Code:    code,
		Message: msg,
		ERR:     err,
	}
	if len(status) > 0 {
		res.StatusCode = status[0]
	}
	return res
}

// Wrap400Response 包装错误码为400的响应错误
func Wrap400Response(err error, msg ...string) error {
	m := "请求发生错误"
	if len(msg) > 0 {
		m = msg[0]
	}
	return WrapResponse(err, 400, m, 400)
}

// Wrap500Response 包装错误码为500的响应错误
func Wrap500Response(err error, msg ...string) error {
	m := "服务器发生错误"
	if len(msg) > 0 {
		m = msg[0]
	}
	return WrapResponse(err, 500, m, 500)
}

// NewResponse 创建响应错误
func NewResponse(code int, msg string, status ...int) error {
	res := &ResponseError{
		Code:    code,
		Message: msg,
	}
	if len(status) > 0 {
		res.StatusCode = status[0]
	}
	return res
}

// New400Response 创建错误码为400的响应错误
func New400Response(msg string) error {
	return NewResponse(400, msg, 400)
}

// New500Response 创建错误码为500的响应错误
func New500Response(msg string) error {
	return NewResponse(500, msg, 500)
}
