package schema

// LoginParam 登录参数
type LoginParam struct {
	UserName   string `json:"user_name" binding:"required" swaggo:"true,用户名"`
	Password   string `json:"password" binding:"required" swaggo:"true,密码(md5加密)"`
	VerifyID   string `json:"verify_id" binding:"required" swaggo:"true,图形验证码ID"`
	VerifyCode string `json:"verify_code" binding:"required" swaggo:"true,图形验证码"`
}

// UserLoginInfo 用户登录信息
type UserLoginInfo struct {
	UserName  string   `json:"user_name" swaggo:"true,用户名"`
	RealName  string   `json:"real_name" swaggo:"true,真实姓名"`
	RoleNames []string `json:"role_names" swaggo:"true,角色名列表"`
}
