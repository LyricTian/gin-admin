package ctl

// Common API模块
type Common struct {
	DemoAPI  *Demo `inject:""`
	LoginAPI *Login
	UserAPI  *User
	RoleAPI  *Role
	MenuAPI  *Menu
}
