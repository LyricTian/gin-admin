package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	_ "github.com/LyricTian/gin-admin/v10/internal/swagger"
	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/LyricTian/gin-admin/v10/internal/wirex"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"go.uber.org/zap"
)

type RunConfig struct {
	ConfigDir string
	StaticDir string
	Daemon    bool
}

func Run(ctx context.Context, cfg RunConfig) error {
	defer zap.L().Sync()

	cfgDir := cfg.ConfigDir
	staticDir := cfg.StaticDir
	config.MustLoad(filepath.Join(cfgDir, "config.toml"))
	config.C.General.ConfigDir = cfgDir
	config.C.Middleware.Static.Dir = staticDir

	if cfg.Daemon {
		bin, err := filepath.Abs(os.Args[0])
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Failed to get absolute path for command: %s", err.Error()))
			return err
		}

		command := exec.Command(bin, "start", "--config", cfgDir, "--static", staticDir)
		err = command.Start()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Failed to start daemon thread: %s", err.Error()))
			return err
		}

		pid := command.Process.Pid
		_ = os.WriteFile(fmt.Sprintf("%s.lock", config.C.General.AppName), []byte(fmt.Sprintf("%d", pid)), 0666)
		os.Stdout.WriteString(fmt.Sprintf("Service %s daemon thread started with pid %d", config.C.General.AppName, pid))
		os.Exit(0)
	}

	config.C.Print()
	cleanLoggerFn, err := InitLogger(ctx, cfgDir)
	if err != nil {
		return err
	}

	ctx = logging.NewTag(ctx, logging.TagKeyMain)
	logging.Context(ctx).Info("Starting service", zap.String("config", cfgDir),
		zap.String("static", staticDir),
		zap.Int("pid", os.Getpid()),
	)

	return utils.Run(ctx, func(ctx context.Context) (func(), error) {
		cleanStartFn, err := Start(ctx)
		if err != nil {
			return nil, err
		}

		return func() {
			if cleanStartFn != nil {
				cleanStartFn()
			}

			if cleanLoggerFn != nil {
				cleanLoggerFn()
			}
		}, nil
	})
}

// The Start function initializes an injector and starts an HTTP server, with an optional pprof server,
// and returns a cleanup function.
func Start(ctx context.Context) (func(), error) {
	injector, cleanInjectorFn, err := wirex.BuildInjector(ctx)
	if err != nil {
		return nil, err
	}

	if err := injector.M.Init(ctx); err != nil {
		return nil, err
	}

	cleanHTTPServerFn, err := startHTTPServer(ctx, injector)
	if err != nil {
		return cleanInjectorFn, err
	}

	if addr := config.C.General.PprofAddr; addr != "" {
		logging.Context(ctx).Info("Pprof server is listening on " + addr)
		go func() {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				logging.Context(ctx).Error("Failed to listen pprof server", zap.Error(err))
			}
		}()
	}

	return func() {
		if cleanHTTPServerFn != nil {
			cleanHTTPServerFn()
		}

		if cleanInjectorFn != nil {
			cleanInjectorFn()
		}
	}, nil
}
