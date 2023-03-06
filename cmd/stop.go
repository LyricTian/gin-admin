package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func StopCmd() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop server",
		Action: func(c *cli.Context) error {
			lockName := fmt.Sprintf("%s.lock", c.App.Name)
			strb, err := os.ReadFile(lockName)
			if err != nil {
				return err
			}

			command := exec.Command("kill", string(strb))
			err = command.Start()
			if err != nil {
				return err
			}

			err = os.Remove(lockName)
			if err != nil {
				return fmt.Errorf("Can't remove %s.lock. %s", c.App.Name, err.Error())
			}

			fmt.Printf("Service %s stopped \n", c.App.Name)
			return nil
		},
	}
}
