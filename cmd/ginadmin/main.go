package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LyricTian/gin-admin/v9/internal"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/urfave/cli/v2"

	_ "github.com/LyricTian/gin-admin/v9/internal/swagger"
)

// Usage: go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "9.0.0"

// @title ginadmin
// @version 9.0.0
// @description RBAC scaffolding based on GIN + GORM + CASBIN + WIRE.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
// @contact.name LyricTian
// @contact.email tiannianshou@gmail.com
func main() {
	ctx := logger.NewTagContext(context.Background(), "__main__")

	app := cli.NewApp()
	app.Name = "ginadmin"
	app.Version = VERSION
	app.Usage = "RBAC scaffolding based on GIN + GORM + CASBIN + WIRE."
	app.Commands = []*cli.Command{
		newServerCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.WithContext(ctx).Errorf(err.Error())
	}
}

func newServerCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Run http server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config_dir",
				Usage:    "Config directory (config.toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "www",
				Usage: "Static site directory",
			},
		},
		Action: func(c *cli.Context) error {
			return run(ctx, func() (func(), error) {
				return internal.Start(ctx,
					internal.SetConfigDir(c.String("config_dir")),
					internal.SetWWWDir(c.String("www")),
					internal.SetVersion(VERSION))
			})
		},
	}
}

func run(ctx context.Context, handler func() (func(), error)) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFn, err := handler()
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("Received signal[%s]", sig.String())

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}
	cleanFn()
	logger.WithContext(ctx).Infof("Server exit, bye...")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
