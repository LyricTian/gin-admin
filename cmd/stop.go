package cmd

import (
	"github.com/LyricTian/gin-admin/v10/internal/bootstrap"
	"github.com/urfave/cli/v2"
)

func StopCmd() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop server",
		Action: func(c *cli.Context) error {
			if err := bootstrap.StopDaemon(); err != nil {
				panic(err)
			}
			return nil
		},
	}
}
