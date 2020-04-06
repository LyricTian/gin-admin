package initialize

import (
	"errors"
	"sync"

	"github.com/LyricTian/gin-admin/internal/app/config"
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

// InitInjector 初始化注入器
func InitInjector() (*Injector, func(), error) {
	var (
		cleanFunc func()
		err       error
	)
	once.Do(func() {
		switch {
		case config.C.Storage.IsGorm():
			I, cleanFunc, err = BuildGormInjector()
		case config.C.Storage.IsMongo():
			I, cleanFunc, err = BuildMongoInjector()
		default:
			err = errors.New("Unknown storage")
		}
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
