package inject

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

// Inject the initialized structure
type Injector struct {
	Engine *gin.Engine
}
