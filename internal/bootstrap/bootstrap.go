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
	"github.com/LyricTian/gin-admin/v10/internal/wirex"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"go.uber.org/zap"
)

type RunConfig struct {
	ConfigDir  string
	ConfigFile string
	StaticDir  string
	Daemon     bool
}

func Run(ctx context.Context, cfg RunConfig) error {
	defer zap.L().Sync()

	cfgDir := cfg.ConfigDir
	staticDir := cfg.StaticDir
	config.MustLoad(filepath.Join(cfgDir, cfg.ConfigFile))
	config.C.General.ConfigDir = cfgDir
	config.C.Middleware.Static.Dir = staticDir

	if cfg.Daemon {
		bin, err := filepath.Abs(os.Args[0])
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Failed to get absolute path for command: %s", err.Error()))
			return err
		}

		command := exec.Command(bin, "start", "--configdir", cfgDir, "--config", cfg.ConfigFile, "--staticdir", staticDir)
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

	return util.Run(ctx, func(ctx context.Context) (func(), error) {
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
		if err := injector.M.Release(ctx); err != nil {
			logging.Context(ctx).Error("Failed to release mods", zap.Error(err))
		}

		if cleanHTTPServerFn != nil {
			cleanHTTPServerFn()
		}

		if cleanInjectorFn != nil {
			cleanInjectorFn()
		}
	}, nil
}

// The function stops a daemon by reading a lock file, killing the process ID found in the file, and
// removing the lock file.
func StopDaemon() error {
	appName := config.C.General.AppName
	lockName := fmt.Sprintf("%s.lock", appName)
	pid, err := os.ReadFile(lockName)
	if err != nil {
		return err
	}

	command := exec.Command("kill", string(pid))
	err = command.Start()
	if err != nil {
		return err
	}

	err = os.Remove(lockName)
	if err != nil {
		return fmt.Errorf("Can't remove %s.lock. %s", appName, err.Error())
	}

	fmt.Printf("Service %s stopped \n", appName)
	return nil
}