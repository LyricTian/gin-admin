package api

import "github.com/google/wire"

var APISet = wire.NewSet(
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
) // end
