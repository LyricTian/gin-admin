package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// This function creates a CLI command that prints the version number.
func VersionCmd(v string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show version",
		Action: func(_ *cli.Context) error {
			fmt.Println(v)
			return nil
		},
	}
}
