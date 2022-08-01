package main

import (
	"fmt"
	"os"

	"github.com/LyricTian/gin-admin/v9/cmd"

	"github.com/urfave/cli/v2"
)

// @title ginadmin
// @version 9.0.0
// @description A simple, modular, high-performance RBAC development framework built on golang.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
// @contact.name LyricTian
// @contact.email tiannianshou@gmail.com
func main() {
	app := cli.NewApp()
	app.Name = "ginadmin"
	app.Version = cmd.VERSION
	app.Usage = "A simple, modular, high-performance RBAC development framework built on golang."
	app.Commands = []*cli.Command{
		cmd.StartCmd,
		cmd.VersionCmd,
		cmd.StopCmd,
	}
	err := app.Run(os.Args)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Failed to run app: %v \n", err))
	}
}
