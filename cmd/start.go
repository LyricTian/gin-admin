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
func StartCmd(ver string) *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-dir",
				Aliases:     []string{"d"},
				Usage:       "Configuration files directory",
				DefaultText: "configs",
				Value:       "configs",
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Runtime configuration files or directory (relative to config-dir, multiple separated by commas, such as: logging.yaml,server.toml)",
				DefaultText: "dev",
				Value:       "dev",
			},
			&cli.StringFlag{
				Name:    "static-dir",
				Aliases: []string{"s"},
				Usage:   "Static files directory",
			},
			&cli.BoolFlag{
				Name:  "daemon",
				Usage: "Run as a daemon",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("daemon") {
				bin, err := filepath.Abs(os.Args[0])
				if err != nil {
					fmt.Printf("Failed to get absolute path for command: %s \n", err.Error())
					return err
				}

				if err := c.Set("daemon", "false"); err != nil {
					return err
				}

				args := []string{"start"}
				args = append(args, "-d", c.String("config-dir"))
				args = append(args, "-c", c.String("config"))
				args = append(args, "-s", c.String("static-dir"))
				fmt.Printf("Execute command: %s %s\n", bin, strings.Join(args, " "))
				command := exec.Command(bin, args...)
				err = command.Start()
				if err != nil {
					fmt.Printf("Failed to start daemon thread: %s \n", err.Error())
					return err
				}

				pid := command.Process.Pid
				_ = os.WriteFile(fmt.Sprintf("%s.lock", c.App.Name), []byte(fmt.Sprintf("%d", pid)), 0666)
				fmt.Printf("Service %s daemon thread started with pid %d \n", config.C.General.AppName, pid)
				os.Exit(0)
			}

			err := bootstrap.Run(context.Background(), bootstrap.RunConfig{
				ConfigDir: c.String("config-dir"),
				Config:    c.String("config"),
				StaticDir: c.String("static-dir"),
				Version:   ver,
			})
			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}
