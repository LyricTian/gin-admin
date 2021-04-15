package service

import "github.com/google/wire"

// ServiceSet model注入
var ServiceSet = wire.NewSet(
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
