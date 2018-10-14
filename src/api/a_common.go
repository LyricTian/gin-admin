package api

// Common API模块
type Common struct {
	DemoAPI *Demo `inject:""`
	MenuAPI *Menu `inject:""`
}
