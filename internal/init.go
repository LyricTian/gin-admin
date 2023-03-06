package internal

import (
	"context"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/internal/library/utilx"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/LyricTian/gin-admin/v10/pkg/x/gormx"
)

func InitLogger(ctx context.Context, cfgDir string) (func(), error) {
	cfg, err := logging.LoadConfigFromToml(filepath.Join(cfgDir, "logging.toml"))
	if err != nil {
		return nil, err
	}

	return logging.InitWithConfig(ctx, cfg, initLoggerHook)
}

func initLoggerHook(ctx context.Context, cfg *logging.HookConfig) (*logging.Hook, error) {
	extra := cfg.Extra
	if extra == nil {
		extra = make(map[string]string)
	}
	extra["appname"] = config.C.General.AppName

	switch cfg.Type {
	case "gorm":
		db, err := gormx.New(gormx.Config{
			Debug:        utilx.DefaultStrToBool(cfg.Options["Debug"], false),
			DBType:       utilx.DefaultStr(cfg.Options["DBType"], "sqlite3"),
			DSN:          utilx.DefaultStr(cfg.Options["DSN"], "data/log.db"),
			MaxLifetime:  utilx.DefaultStrToInt(cfg.Options["MaxLifetime"], 86400),
			MaxIdleTime:  utilx.DefaultStrToInt(cfg.Options["MaxIdleTime"], 3600),
			MaxOpenConns: utilx.DefaultStrToInt(cfg.Options["MaxOpenConns"], 8),
			MaxIdleConns: utilx.DefaultStrToInt(cfg.Options["MaxIdleConns"], 4),
			TablePrefix:  utilx.DefaultStr(cfg.Options["TablePrefix"], ""),
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
