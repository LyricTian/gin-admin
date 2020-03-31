package inject

import (
	"sync"

	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var (
	// I 全局注入器
	I    *Injector
	once sync.Once
)

// InitializeInjector 初始化注入器
func InitializeInjector() (*Injector, func(), error) {
	var (
		cleanFunc func()
		err       error
	)
	once.Do(func() {
		I, cleanFunc, err = BuildInjector()
	})
	return I, cleanFunc, err
}

// InjectorSet 注入Injector
var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

// Injector 全局注入器
type Injector struct {
	Engine         *gin.Engine
	Auth           auth.Auther
	CasbinEnforcer *casbin.SyncedEnforcer
	Menu           *Menu
}
