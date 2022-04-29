package consts

const (
	// common.http
	ErrBadRequestID          = "com.http.bad_request"
	ErrUnauthorizedID        = "com.http.unauthorized"
	ErrForbiddenID           = "com.http.forbidden"
	ErrNotFoundID            = "com.http.not_found"
	ErrMethodNotAllowedID    = "com.http.method_not_allowed"
	ErrTooManyRequestsID     = "com.http.too_many_requests"
	ErrRequestEntityTooLarge = "com.http.request_entity_too_large"
	ErrInternalServerErrorID = "com.http.internal_server_error"

	// common.perm
	ErrInvalidTokenID     = "com.perm.invalid_token"
	ErrNoDataPermissionID = "com.perm.unauthorized_data"

	// common.verify
	ErrInvalidCaptchaID = "com.verify.invalid_captcha"

	// common.hash
	ErrHashPasswordID = "com.hash.password"

	// user
	ErrUserNotFoundID = "user.not_found"
	ErrUserExists     = "user.exists"
)

// define logger tag
const (
	LogLoginTag    = "__login__"
	LogRegisterTag = "__register__"
	LogLogoutTag   = "__logout__"
	LogPasswordTag = "__password__"
)
