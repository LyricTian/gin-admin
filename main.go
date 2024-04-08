package main

import (
	"os"

	"github.com/LyricTian/gin-admin/v10/cmd"
	"github.com/urfave/cli/v2"
)

// VERSION Usage: go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "v10.0.2"

// @title ginadmin
// @version v10.0.2
// @description A lightweight, flexible, elegant and full-featured RBAC scaffolding based on GIN + GORM 2.0 + Casbin 2.0 + Wire DI.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
func main() {
	app := cli.NewApp()
	app.Name = "ginadmin"
	app.Version = VERSION
	app.Usage = "A lightweight, flexible, elegant and full-featured RBAC scaffolding based on GIN + GORM 2.0 + Casbin 2.0 + Wire DI."
	app.Commands = []*cli.Command{
		cmd.StartCmd(),
		cmd.StopCmd(),
		cmd.VersionCmd(VERSION),
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
