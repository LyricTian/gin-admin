package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v10/internal"
	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/internal/library/utilx"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func StartCmd() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "configdir",
				Aliases: []string{"c"},
				Usage:   "Configuration directory (config.toml)",
				Value:   "configs",
			},
			&cli.StringFlag{
				Name:    "staticdir",
				Aliases: []string{"s"},
				Usage:   "Static site directory",
			},
			&cli.BoolFlag{
				Name:    "daemon",
				Aliases: []string{"d"},
				Usage:   "Run as a daemon",
			},
		},
		Action: func(c *cli.Context) error {
			defer func() {
				_ = zap.L().Sync()
			}()

			ctx := logging.NewTag(context.Background(), logging.TagKeyMain)
			return utilx.Run(ctx, func(ctx context.Context) (func(), error) {
				cfgDir := c.String("configdir")
				staticDir := c.String("staticdir")
				daemon := c.Bool("daemon")

				if daemon {
					bin, err := filepath.Abs(os.Args[0])
					if err != nil {
						logging.Context(ctx).Error("Failed to get absolute path for command", zap.Error(err))
						return nil, err
					}

					command := exec.Command(bin, "start", "--configdir", cfgDir, "--staticdir", staticDir)
					err = command.Start()
					if err != nil {
						logging.Context(ctx).Error("Failed to start daemon thread", zap.Error(err))
						return nil, err
					}
					_ = os.WriteFile(fmt.Sprintf("%s.lock", c.App.Name), []byte(fmt.Sprintf("%d", command.Process.Pid)), 0666)
					os.Exit(0)
				}

				// Load configurations
				config.MustLoad(filepath.Join(cfgDir, "config.toml"))
				config.C.General.ConfigDir = cfgDir
				config.C.Middleware.Static.Dir = staticDir
				config.C.Print()

				// Init logger
				loggerClean, err := internal.InitLogger(ctx, cfgDir)
				if err != nil {
					return nil, err
				}

				logging.Context(ctx).Info("Starting server",
					zap.String("configdir", cfgDir),
					zap.String("staticdir", staticDir),
					zap.Bool("daemon", daemon),
					zap.Int("pid", os.Getpid()),
				)

				// Start server
				startClean, err := internal.Start(ctx)
				if err != nil {
					return nil, err
				}

				return func() {
					if startClean != nil {
						startClean()
					}
					if loggerClean != nil {
						loggerClean()
					}
				}, nil
			})
		},
	}
}
