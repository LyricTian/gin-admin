package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v9/internal"
	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/LyricTian/gin-admin/v9/pkg/x/gormx"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var StartCmd = &cli.Command{
	Name:  "start",
	Usage: "Start server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "configdir",
			Usage:    "Configuration directory (logger.toml/config.toml)",
			Required: false,
			Value:    "configs",
		},
		&cli.StringFlag{
			Name:  "staticdir",
			Usage: "Static site directory",
		},
		&cli.BoolFlag{
			Name:    "deamon",
			Aliases: []string{"d"},
			Usage:   "Run as a deamon",
		},
	},
	Action: func(c *cli.Context) error {
		defer func() {
			err := zap.L().Sync()
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Failed to sync zap logger: %v \n", err))
			}
		}()

		ctx := logger.NewTag(context.Background(), "start")
		return utilx.Run(ctx, func() (func(), error) {
			cfgDir := c.String("configdir")
			if cfgDir == "" {
				cfgDir = "configs"
			}

			// Initialize the logger first
			cleanLogger, err := logger.InitWithConfig(filepath.Join(cfgDir, "logger.toml"), HandleLoggerHook)
			if err != nil {
				return nil, err
			}

			daemon := c.Bool("deamon")
			staticDir := c.String("staticdir")
			logger.Context(ctx).Info("Starting server",
				zap.String("configdir", cfgDir),
				zap.String("staticdir", staticDir),
				zap.Int("pid", os.Getpid()),
				zap.Bool("deamon", daemon),
			)

			if daemon {
				bin, err := filepath.Abs(os.Args[0])
				if err != nil {
					logger.Context(ctx).Error("Failed to get absolute path", zap.Error(err))
					return nil, err
				}

				command := exec.Command(bin, "start", "--configdir", cfgDir, "--staticdir", c.String("staticdir"))
				err = command.Start()
				if err != nil {
					logger.Context(ctx).Error("Failed to start deamon", zap.Error(err))
					return nil, err
				}
				_ = ioutil.WriteFile(fmt.Sprintf("%s.lock", c.App.Name), []byte(fmt.Sprintf("%d", command.Process.Pid)), 0666)
				os.Exit(0)
			}

			// Load the configuration
			config.MustLoad(filepath.Join(cfgDir, "config.toml"))
			config.C.General.ConfigDir = cfgDir
			config.C.Middleware.Static.Dir = staticDir
			if !config.C.General.DisablePrintConfig {
				logger.Context(ctx).Info("Load configuration", zap.String("config", config.C.String()))
			}

			// Initialize the server
			cleanInit, err := internal.Init(ctx)
			if err != nil {
				return nil, err
			}

			return func() {
				if cleanInit != nil {
					cleanInit()
				}
				if cleanLogger != nil {
					cleanLogger()
				}
			}, nil
		})
	},
}

func HandleLoggerHook(md *toml.MetaData, cfg *logger.HookConfig) (*logger.Hook, error) {
	switch cfg.Type {
	case "gorm":
		var gormxCfg gormx.Config
		err := md.PrimitiveDecode(*cfg.Options, &gormxCfg)
		if err != nil {
			return nil, err
		}

		db, err := gormx.New(gormxCfg)
		if err != nil {
			return nil, err
		}

		hook := logger.NewHook(logger.NewGormHook(db),
			logger.SetHookExtra(cfg.Extra),
			logger.SetHookMaxJobs(cfg.MaxBuffer),
			logger.SetHookMaxWorkers(cfg.MaxThread))
		return hook, nil
	default:
		return nil, nil
	}
}
