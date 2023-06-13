package bootstrap

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/pkg/gormx"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/spf13/cast"
)

func initLoggerHook(_ context.Context, cfg *logging.HookConfig) (*logging.Hook, error) {
	extra := cfg.Extra
	if extra == nil {
		extra = make(map[string]string)
	}
	extra["appname"] = config.C.General.AppName

	switch cfg.Type {
	case "gorm":
		db, err := gormx.New(gormx.Config{
			Debug:        cast.ToBool(cfg.Options["Debug"]),
			DBType:       cast.ToString(cfg.Options["DBType"]),
			DSN:          cast.ToString(cfg.Options["DSN"]),
			MaxLifetime:  cast.ToInt(cfg.Options["MaxLifetime"]),
			MaxIdleTime:  cast.ToInt(cfg.Options["MaxIdleTime"]),
			MaxOpenConns: cast.ToInt(cfg.Options["MaxOpenConns"]),
			MaxIdleConns: cast.ToInt(cfg.Options["MaxIdleConns"]),
			TablePrefix:  cast.ToString(cfg.Options["TablePrefix"]),
		})
		if err != nil {
			return nil, err
		}

		hook := logging.NewHook(logging.NewGormHook(db),
			logging.SetHookExtra(cfg.Extra),
			logging.SetHookMaxJobs(cfg.MaxBuffer),
			logging.SetHookMaxWorkers(cfg.MaxThread))
		return hook, nil
	default:
		return nil, nil
	}
}
