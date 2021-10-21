package mock

import "github.com/google/wire"

var MockSet = wire.NewSet(
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
) // end
