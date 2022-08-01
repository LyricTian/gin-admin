package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/LyricTian/gin-admin/v9/pkg/logger"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// Usage: go build -ldflags "-X github.com/LyricTian/gin-admin/v9/cmd.VERSION=x.x.x"
var VERSION = "v9.0.0"

var VersionCmd = &cli.Command{
	Name:  "version",
	Usage: "Show version",
	Action: func(c *cli.Context) error {
		fmt.Println(VERSION)
		return nil
	},
}

var StopCmd = &cli.Command{
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
			return fmt.Errorf("Can't remove signaling.lock. %s", err.Error())
		}

		fmt.Printf(" %s stopped \n", c.App.Name)
		return nil
	},
}

func Run(ctx context.Context, handler func() (func(), error)) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFn, err := handler()
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.Context(ctx).Info("Received signal", zap.String("signal", sig.String()))

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFn()
	logger.Context(ctx).Info("Server exit, bye...")
	time.Sleep(time.Millisecond * 100)
	os.Exit(state)
	return nil
}
