package bootstrap

import (
	"context"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/internal/library/utils"
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
			Debug:        utils.DefaultStrToBool(cfg.Options["Debug"], false),
			DBType:       utils.DefaultStr(cfg.Options["DBType"], "sqlite3"),
			DSN:          utils.DefaultStr(cfg.Options["DSN"], "data/log.db"),
			MaxLifetime:  utils.DefaultStrToInt(cfg.Options["MaxLifetime"], 86400),
			MaxIdleTime:  utils.DefaultStrToInt(cfg.Options["MaxIdleTime"], 3600),
			MaxOpenConns: utils.DefaultStrToInt(cfg.Options["MaxOpenConns"], 8),
			MaxIdleConns: utils.DefaultStrToInt(cfg.Options["MaxIdleConns"], 4),
			TablePrefix:  utils.DefaultStr(cfg.Options["TablePrefix"], ""),
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
