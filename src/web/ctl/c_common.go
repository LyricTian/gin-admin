package ctl

// Common API模块
type Common struct {
	DemoAPI  *Demo  `inject:""`
	LoginAPI *Login `inject:""`
	UserAPI  *User  `inject:""`
	RoleAPI  *Role  `inject:""`
	MenuAPI  *Menu  `inject:""`
}
