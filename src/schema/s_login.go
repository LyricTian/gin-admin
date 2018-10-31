package schema

// LoginParam 登录参数
type LoginParam struct {
	UserName string `json:"user_name" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"`  // 密码(md5加密)
}
