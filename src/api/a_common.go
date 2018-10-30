package api

// Common API模块
type Common struct {
	UserAPI *User `inject:""`
	RoleAPI *Role `inject:""`
	DemoAPI *Demo `inject:""`
	MenuAPI *Menu `inject:""`
}
