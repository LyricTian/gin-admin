package service

import "github.com/google/wire"

// ServiceSet bll注入
var ServiceSet = wire.NewSet(
	DemoSet,
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
)
