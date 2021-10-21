package service

import (
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	MenuSet,
	RoleSet,
	UserSet,
	LoginSet,
) // end
