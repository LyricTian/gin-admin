package bll

import "github.com/google/wire"

// AllSet bll注入
var AllSet = wire.NewSet(
	DemoSet,
	LoginSet,
	MenuSet,
	RoleSet,
	TransSet,
	UserSet,
)
