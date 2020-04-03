package main

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/LyricTian/gin-admin/internal/app"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/urfave/cli/v2"
)

// VERSION 版本号，
// 可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "6.0.0"

func main() {
	logger.SetVersion(VERSION)
	logger.SetTraceIDFunc(util.NewTraceID)
	ctx := logger.NewTraceIDContext(context.Background(), util.NewTraceID())

	app := cli.NewApp()
	app.Name = "gin-admin"
	app.Version = VERSION
	app.Usage = "RBAC scaffolding based on Gin + Gorm + Casbin + Wire."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf(ctx, err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "启动HTTP服务",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "配置文件(.json,.yaml,.toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "model",
				Aliases:  []string{"m"},
				Usage:    "casbin的访问控制模型(.conf)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "menu",
				Usage: "初始化菜单数据配置文件(.yaml)",
			},
			&cli.StringFlag{
				Name:  "www",
				Usage: "静态站点目录",
			},
		},
		Action: func(c *cli.Context) error {
			var state int32 = 1
			sc := make(chan os.Signal, 1)
			signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			cleanFunc := app.Init(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetModelFile(c.String("model")),
				app.SetWWWDir(c.String("www")),
				app.SetMenuFile(c.String("menu_data")),
				app.SetVersion(VERSION))

		EXIT:
			for {
				sig := <-sc
				logger.Printf(ctx, "获取到信号[%s]", sig.String())
				switch sig {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					atomic.CompareAndSwapInt32(&state, 1, 0)
					break EXIT
				case syscall.SIGHUP:
				default:
					break EXIT
				}
			}

			cleanFunc()
			logger.Printf(ctx, "服务退出")
			time.Sleep(time.Second)
			os.Exit(int(atomic.LoadInt32(&state)))
			return nil
		},
	}
}
