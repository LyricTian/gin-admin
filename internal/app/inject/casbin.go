package inject

import (
	"time"

	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/module/adapter"
	"github.com/casbin/casbin/v2"
)

// InitCasbin 初始化casbin
func InitCasbin(adapter *adapter.CasbinAdapter) (*casbin.SyncedEnforcer, func(), error) {
	cfg := config.C.Casbin
	if cfg.Model == "" {
		return new(casbin.SyncedEnforcer), nil, nil
	}

	e, err := casbin.NewSyncedEnforcer(cfg.Model, cfg.Debug)
	if err != nil {
		return nil, nil, err
	}

	err = e.InitWithModelAndAdapter(e.GetModel(), adapter)
	if err != nil {
		return nil, nil, err
	}
	e.EnableEnforce(cfg.Enable)

	var cleanFunc func()
	if cfg.AutoLoad {
		e.StartAutoLoadPolicy(time.Duration(cfg.AutoLoadInternal) * time.Second)
		cleanFunc = func() {
			e.StopAutoLoadPolicy()
		}
	}

	return e, cleanFunc, nil
}
