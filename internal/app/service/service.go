package service

import (
	"github.com/google/wire"
)

// ServiceSet bll注入
var ServiceSet = wire.NewSet(
	MenuSet,
	RoleSet,
	UserSet,
	LoginSet,
) // end
