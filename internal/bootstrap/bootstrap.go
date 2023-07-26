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

// RunConfig defines the config for run command.
type RunConfig struct {
	ConfigDir string // Configurations directory
	Config    string // Directory or files (multiple separated by commas)
	StaticDir string // Static files directory
	Version   string
}

// The Run function initializes and starts a service with configuration and logging, and handles
// cleanup upon exit.
func Run(ctx context.Context, runCfg RunConfig) error {
	defer func() {
		if err := zap.L().Sync(); err != nil {
			fmt.Printf("Failed to sync zap logger: %s \n", err.Error())
		}
	}()

	cfgDir := runCfg.ConfigDir
	staticDir := runCfg.StaticDir
	config.MustLoad(cfgDir, strings.Split(runCfg.Config, ",")...)
	config.C.General.ConfigDir = cfgDir
	config.C.Middleware.Static.Dir = staticDir
	config.C.General.Version = runCfg.Version
	config.C.Print()

	cleanLoggerFn, err := logging.InitWithConfig(ctx, &config.C.Logger, initLoggerHook)
	if err != nil {
		return err
	}
	ctx = logging.NewTag(ctx, logging.TagKeyMain)

	logging.Context(ctx).Info("Starting service ...",
		zap.String("version", runCfg.Version),
		zap.Int("pid", os.Getpid()),
		zap.String("config_dir", cfgDir),
		zap.String("config", runCfg.Config),
		zap.String("static_dir", staticDir),
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

	// Initialize injector
	injector, cleanInjectorFn, err := wirex.BuildInjector(ctx)
	if err != nil {
		return err
	}
	if err := injector.M.Init(ctx); err != nil {
		return err
	}

	return util.Run(ctx, func(ctx context.Context) (func(), error) {
		// Start HTTP server
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
