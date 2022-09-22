package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/LyricTian/gin-admin/v9/cmd"
	"github.com/urfave/cli/v2"
)

// Usage: go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "v9.0.1"

// @title Bobber DevOps
// @version v9.0.1
// @description A DevOps platform based on golang for service monitoring, configuration center, log search, etc.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
// @contact.name LyricTian
// @contact.email lyric.tian@bobcatminer.com
func main() {
	app := cli.NewApp()
	app.Name = "devops"
	app.Version = VERSION
	app.Usage = "A DevOps platform based on golang for service monitoring, configuration center, log search, etc."
	app.Commands = []*cli.Command{
		{
			Name:  "version",
			Usage: "Show version",
			Action: func(c *cli.Context) error {
				fmt.Println(VERSION)
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Stop server",
			Action: func(c *cli.Context) error {
				lockFileName := fmt.Sprintf("%s.lock", c.App.Name)
				strb, err := ioutil.ReadFile(lockFileName)
				if err != nil {
					return err
				}
				command := exec.Command("kill", string(strb))
				err = command.Start()
				if err != nil {
					return err
				}

				err = os.Remove(lockFileName)
				if err != nil {
					return fmt.Errorf("Can't remove %s.lock. %s", c.App.Name, err.Error())
				}

				fmt.Printf("Service %s stopped \n", c.App.Name)
				return nil
			},
		},
		cmd.StartCmd,
		cmd.GenerateCmd,
	}
	err := app.Run(os.Args)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Failed to run app: %v \n", err))
	}
}
