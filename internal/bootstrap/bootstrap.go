package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	_ "github.com/LyricTian/gin-admin/v10/internal/swagger"
	"github.com/LyricTian/gin-admin/v10/internal/utility/prom"
	"github.com/LyricTian/gin-admin/v10/internal/wirex"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"go.uber.org/zap"
)

// RunConfig defines the config for run command.
type RunConfig struct {
	WorkDir   string // Working directory
	Configs   string // Directory or files (multiple separated by commas)
	StaticDir string // Static files directory
}

// The Run function initializes and starts a service with configuration and logging, and handles
// cleanup upon exit.
func Run(ctx context.Context, runCfg RunConfig) error {
	defer func() {
		if err := zap.L().Sync(); err != nil {
			fmt.Printf("failed to sync zap logger: %s \n", err.Error())
		}
	}()

	// Load configuration.
	workDir := runCfg.WorkDir
	staticDir := runCfg.StaticDir
	config.MustLoad(workDir, strings.Split(runCfg.Configs, ",")...)
	config.C.General.WorkDir = workDir
	config.C.Middleware.Static.Dir = staticDir
	config.C.Print()
	config.C.PreLoad()

	// Initialize logger.
	cleanLoggerFn, err := logging.InitWithConfig(ctx, &config.C.Logger, initLoggerHook)
	if err != nil {
		return err
	}
	ctx = logging.NewTag(ctx, logging.TagKeyMain)

	logging.Context(ctx).Info("starting service ...",
		zap.String("version", config.C.General.Version),
		zap.Int("pid", os.Getpid()),
		zap.String("workdir", workDir),
		zap.String("config", runCfg.Configs),
		zap.String("static", staticDir),
	)

	// Start pprof server.
	if addr := config.C.General.PprofAddr; addr != "" {
		logging.Context(ctx).Info("pprof server is listening on " + addr)
		go func() {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				logging.Context(ctx).Error("failed to listen pprof server", zap.Error(err))
			}
		}()
	}

	// Build injector.
	injector, cleanInjectorFn, err := wirex.BuildInjector(ctx)
	if err != nil {
		return err
	}

	if err := injector.M.Init(ctx); err != nil {
		return err
	}

	// Initialize global prometheus metrics.
	prom.Init()

	return util.Run(ctx, func(ctx context.Context) (func(), error) {
		cleanHTTPServerFn, err := startHTTPServer(ctx, injector)
		if err != nil {
			return cleanInjectorFn, err
		}

		return func() {
			if err := injector.M.Release(ctx); err != nil {
				logging.Context(ctx).Error("failed to release injector", zap.Error(err))
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
