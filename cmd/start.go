package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v10/internal/bootstrap"
	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/internal/utils"
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
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Configuration directory",
				Value:   "configs",
			},
			&cli.StringFlag{
				Name:    "static",
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

			cfgDir := c.String("config")
			staticDir := c.String("static")

			if c.Bool("daemon") {
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
				_ = os.WriteFile(fmt.Sprintf("%s.lock", c.App.Name), []byte(fmt.Sprintf("%d", pid)), 0666)
				os.Stdout.WriteString(fmt.Sprintf("Daemon thread started with pid %d", pid))
				os.Exit(0)
			}

			ctx := logging.NewTag(context.Background(), logging.TagKeyMain)
			return utils.Run(ctx, func(ctx context.Context) (func(), error) {
				config.MustLoad(filepath.Join(cfgDir, "config.toml"))
				config.C.General.ConfigDir = cfgDir
				config.C.Middleware.Static.Dir = staticDir
				config.C.Print()

				cleanLoggerFn, err := bootstrap.InitLogger(ctx, cfgDir)
				if err != nil {
					return nil, err
				}

				logging.Context(ctx).Info("Starting server", zap.String("config", cfgDir),
					zap.String("static", staticDir),
					zap.Int("pid", os.Getpid()),
				)

				cleanStartFn, err := bootstrap.Start(ctx)
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
		},
	}
}
