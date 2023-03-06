package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func VersionCmd(v string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show version",
		Action: func(c *cli.Context) error {
			fmt.Println(v)
			return nil
		},
	}
}
