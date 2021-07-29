package api

import "github.com/google/wire"

// APISet 注入api
var APISet = wire.NewSet(
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
)
