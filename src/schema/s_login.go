package schema

// LoginParam 登录参数
type LoginParam struct {
	UserName string `json:"user_name" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"`  // 密码(md5加密)
}

// LoginInfo 用户登录信息
type LoginInfo struct {
	UserName  string   `json:"user_name"`  // 用户名
	RealName  string   `json:"real_name"`  // 真实姓名
	RoleNames []string `json:"role_names"` // 真实姓名
}
