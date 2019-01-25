package schema

// LoginParam 登录参数
type LoginParam struct {
	UserName   string `json:"user_name" binding:"required" swaggo:"true,用户名"`
	Password   string `json:"password" binding:"required" swaggo:"true,密码(md5加密)"`
	VerifyID   string `json:"verify_id" binding:"required" swaggo:"true,验证码ID"`
	VerifyCode string `json:"verify_code" binding:"required" swaggo:"true,验证码"`
}

// UserLoginInfo 用户登录信息
type UserLoginInfo struct {
	UserName  string   `json:"user_name" swaggo:"true,用户名"`
	RealName  string   `json:"real_name" swaggo:"true,真实姓名"`
	RoleNames []string `json:"role_names" swaggo:"true,角色名列表"`
}

// UpdatePasswordParam 更新个人密码请求参数
type UpdatePasswordParam struct {
	OldPassword string `json:"old_password" binding:"required" swaggo:"true,旧密码(md5加密)"`
	NewPassword string `json:"new_password" binding:"required" swaggo:"true,新密码(md5加密)"`
}

// LoginCaptcha 登录验证码
type LoginCaptcha struct {
	CaptchaID string `json:"captcha_id" swaggo:"true,验证码ID"`
}
