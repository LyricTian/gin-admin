package schema

type LoginParam struct {
	UserName    string `json:"user_name" binding:"required"`    // 用户名
	Password    string `json:"password" binding:"required"`     // 密码(md5加密)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // 验证码ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // 验证码
}

type UserLoginInfo struct {
	UserID   uint64 `json:"user_id,string"` // 用户ID
	UserName string `json:"user_name"`      // 用户名
	RealName string `json:"real_name"`      // 真实姓名
	Roles    Roles  `json:"roles"`          // 角色列表
}

type UpdatePasswordParam struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码(md5加密)
	NewPassword string `json:"new_password" binding:"required"` // 新密码(md5加密)
}

type LoginCaptcha struct {
	CaptchaID string `json:"captcha_id"` // 验证码ID
}

type LoginTokenInfo struct {
	AccessToken string `json:"access_token"` // 访问令牌
	TokenType   string `json:"token_type"`   // 令牌类型
	ExpiresAt   int64  `json:"expires_at"`   // 过期时间戳
}
