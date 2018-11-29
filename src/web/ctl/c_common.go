package ctl

// Common API模块
type Common struct {
	LoginAPI *Login `inject:""`
	UserAPI  *User  `inject:""`
	RoleAPI  *Role  `inject:""`
	DemoAPI  *Demo  `inject:""`
	MenuAPI  *Menu  `inject:""`
}
