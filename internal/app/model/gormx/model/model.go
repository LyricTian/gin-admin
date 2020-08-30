package model

import "github.com/google/wire"

// ModelSet model注入
var ModelSet = wire.NewSet(
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
