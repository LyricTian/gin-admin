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

// The function defines a CLI command to start a server with various flags and options, including the
// ability to run as a daemon.
func StartCmd() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "workdir",
				Aliases:     []string{"d"},
				Usage:       "Working directory",
				DefaultText: "configs",
				Value:       "configs",
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Runtime configuration files or directory (relative to workdir, multiple separated by commas)",
				DefaultText: "dev",
				Value:       "dev",
			},
			&cli.StringFlag{
				Name:    "static",
				Aliases: []string{"s"},
				Usage:   "Static files directory",
			},
			&cli.BoolFlag{
				Name:  "daemon",
				Usage: "Run as a daemon",
			},
		},
		Action: func(c *cli.Context) error {
			workDir := c.String("workdir")
			staticDir := c.String("static")
			configs := c.String("config")

			if c.Bool("daemon") {
				bin, err := filepath.Abs(os.Args[0])
				if err != nil {
					fmt.Printf("failed to get absolute path for command: %s \n", err.Error())
					return err
				}

				args := []string{"start"}
				args = append(args, "-d", workDir)
				args = append(args, "-c", configs)
				args = append(args, "-s", staticDir)
				fmt.Printf("execute command: %s %s \n", bin, strings.Join(args, " "))
				command := exec.Command(bin, args...)
				err = command.Start()
				if err != nil {
					fmt.Printf("failed to start daemon thread: %s \n", err.Error())
					return err
				}

				pid := command.Process.Pid
				_ = os.WriteFile(fmt.Sprintf("%s.lock", c.App.Name), []byte(fmt.Sprintf("%d", pid)), 0666)
				fmt.Printf("service %s daemon thread started with pid %d \n", config.C.General.AppName, pid)
				os.Exit(0)
			}

			err := bootstrap.Run(context.Background(), bootstrap.RunConfig{
				WorkDir:   workDir,
				Configs:   configs,
				StaticDir: staticDir,
			})
			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}
