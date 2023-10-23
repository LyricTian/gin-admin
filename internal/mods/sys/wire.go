package sys

import (
	"github.com/LyricTian/gin-admin/v10/internal/mods/sys/api"
	"github.com/LyricTian/gin-admin/v10/internal/mods/sys/biz"
	"github.com/LyricTian/gin-admin/v10/internal/mods/sys/dal"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(SYS), "*"),
	wire.Struct(new(dal.Logger), "*"),
	wire.Struct(new(biz.Logger), "*"),
	wire.Struct(new(api.Logger), "*"),
)
