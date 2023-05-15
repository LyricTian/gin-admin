package bootstrap

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/pkg/gormx"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
)

func initLoggerHook(ctx context.Context, cfg *logging.HookConfig) (*logging.Hook, error) {
	extra := cfg.Extra
	if extra == nil {
		extra = make(map[string]string)
	}
	extra["appname"] = config.C.General.AppName

	switch cfg.Type {
	case "gorm":
		db, err := gormx.New(gormx.Config{
			Debug:        util.DefaultStrToBool(cfg.Options["Debug"], false),
			DBType:       util.DefaultStr(cfg.Options["DBType"], "sqlite3"),
			DSN:          util.DefaultStr(cfg.Options["DSN"], "data/log.db"),
			MaxLifetime:  util.DefaultStrToInt(cfg.Options["MaxLifetime"], 86400),
			MaxIdleTime:  util.DefaultStrToInt(cfg.Options["MaxIdleTime"], 3600),
			MaxOpenConns: util.DefaultStrToInt(cfg.Options["MaxOpenConns"], 8),
			MaxIdleConns: util.DefaultStrToInt(cfg.Options["MaxIdleConns"], 4),
			TablePrefix:  util.DefaultStr(cfg.Options["TablePrefix"], ""),
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
