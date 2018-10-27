package api

// Common API模块
type Common struct {
	RoleAPI *Role `inject:""`
	DemoAPI *Demo `inject:""`
	MenuAPI *Menu `inject:""`
}
