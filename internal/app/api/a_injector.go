package api

import "github.com/google/wire"

// AllSet 注入api
var AllSet = wire.NewSet(
	DemoSet,
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
)
