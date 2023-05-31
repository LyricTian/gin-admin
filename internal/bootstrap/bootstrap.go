package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	_ "github.com/LyricTian/gin-admin/v10/internal/swagger"
	"github.com/LyricTian/gin-admin/v10/internal/wirex"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"go.uber.org/zap"
)

type RunConfig struct {
	ConfigDir  string
	ConfigFile string
	StaticDir  string
}

// The Run function initializes and starts a service with configuration and logging.
func Run(ctx context.Context, runCfg RunConfig) error {
	defer func() {
		if err := zap.L().Sync(); err != nil {
			fmt.Printf("Failed to sync zap logger: %s \n", err.Error())
		}
	}()

	cfgDir := runCfg.ConfigDir
	staticDir := runCfg.StaticDir
	config.MustLoad(cfgDir, strings.Split(runCfg.ConfigFile, ",")...)
	config.C.General.ConfigDir = cfgDir
	config.C.Middleware.Static.Dir = staticDir
	config.C.Print()

	cleanLoggerFn, err := logging.InitWithConfig(ctx, &config.C.Logger, initLoggerHook)
	if err != nil {
		return err
	}

	ctx = logging.NewTag(ctx, logging.TagKeyMain)
	logging.Context(ctx).Info("Starting service",
		zap.Any("runConfig", runCfg),
		zap.Int("pid", os.Getpid()),
	)

	if addr := config.C.General.PprofAddr; addr != "" {
		logging.Context(ctx).Info("Pprof server is listening on " + addr)
		go func() {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				logging.Context(ctx).Error("Failed to listen pprof server", zap.Error(err))
			}
		}()
	}

	injector, cleanInjectorFn, err := wirex.BuildInjector(ctx)
	if err != nil {
		return err
	}

	if err := injector.M.Init(ctx); err != nil {
		return err
	}

	return util.Run(ctx, func(ctx context.Context) (func(), error) {
		cleanHTTPServerFn, err := startHTTPServer(ctx, injector)
		if err != nil {
			return cleanInjectorFn, err
		}

		return func() {
			if err := injector.M.Release(ctx); err != nil {
				logging.Context(ctx).Error("Failed to release mods", zap.Error(err))
			}
			if cleanHTTPServerFn != nil {
				cleanHTTPServerFn()
			}
			if cleanInjectorFn != nil {
				cleanInjectorFn()
			}
			if cleanLoggerFn != nil {
				cleanLoggerFn()
			}
		}, nil
	})
}
