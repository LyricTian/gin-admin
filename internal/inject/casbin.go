package inject

import (
	"context"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/casbin/casbin/v2"
)

func InitCasbin(ctx context.Context) (*casbin.Enforcer, func(), error) {
	cfg := config.C.Casbin

	emptyFunc := func() {}
	if !cfg.Enable {
		return new(casbin.Enforcer), emptyFunc, nil
	}

	tempFile := filepath.Join(config.C.ConfigDir, "policy.csv")
	modelFile := filepath.Join(config.C.ConfigDir, "acl_casbin_model.conf")
	e, err := casbin.NewEnforcer(modelFile, tempFile)
	if err != nil {
		return nil, nil, err
	}
	e.EnableLog(cfg.Debug)

	return e, emptyFunc, nil
}
