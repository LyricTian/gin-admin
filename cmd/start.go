package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/LyricTian/gin-admin/v10/internal/bootstrap"
	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/urfave/cli/v2"
)

func StartCmd() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "configdir",
				Usage:       "Configurations directory",
				DefaultText: "configs",
				Value:       "configs",
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Configuration directory or files (multiple separated by commas)",
				DefaultText: "dev",
				Value:       "dev",
			},
			&cli.StringFlag{
				Name:  "staticdir",
				Usage: "Static files directory",
			},
			&cli.BoolFlag{
				Name:    "daemon",
				Aliases: []string{"d"},
				Usage:   "Run as a daemon",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("daemon") {
				bin, err := filepath.Abs(os.Args[0])
				if err != nil {
					os.Stderr.WriteString(fmt.Sprintf("Failed to get absolute path for command: %s", err.Error()))
					return err
				}

				if err := c.Set("daemon", "false"); err != nil {
					return err
				}

				args := []string{"start"}
				args = append(args, "--configdir", c.String("configdir"))
				args = append(args, "--config", c.String("config"))
				args = append(args, "--staticdir", c.String("staticdir"))
				fmt.Printf("Execute command: %s %s\n", bin, strings.Join(args, " "))
				command := exec.Command(bin, args...)
				err = command.Start()
				if err != nil {
					os.Stderr.WriteString(fmt.Sprintf("Failed to start daemon thread: %s", err.Error()))
					return err
				}

				pid := command.Process.Pid
				_ = os.WriteFile(fmt.Sprintf("%s.lock", c.App.Name), []byte(fmt.Sprintf("%d", pid)), 0666)
				os.Stdout.WriteString(fmt.Sprintf("Service %s daemon thread started with pid %d\n", config.C.General.AppName, pid))
				os.Exit(0)
			}

			err := bootstrap.Run(context.Background(), bootstrap.RunConfig{
				ConfigDir:  c.String("configdir"),
				ConfigFile: c.String("config"),
				StaticDir:  c.String("staticdir"),
			})
			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}
