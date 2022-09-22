package utilx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/x/contextx"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"go.uber.org/zap"
)

func IsRootUser(ctx context.Context) bool {
	return contextx.FromUserID(ctx) == config.C.Dictionary.RootUser.ID
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
