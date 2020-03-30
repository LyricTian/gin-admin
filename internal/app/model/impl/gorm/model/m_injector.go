package model

import "github.com/google/wire"

// AllSet model注入
var AllSet = wire.NewSet(
	DemoSet,
	MenuActionResourceSet,
	MenuActionSet,
	MenuSet,
	RoleMenuSet,
	RoleSet,
	TransSet,
	UserRoleSet,
	UserSet,
)
