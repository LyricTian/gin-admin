package bootstrap

import (
	"context"
	"net/http"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	_ "github.com/LyricTian/gin-admin/v10/internal/swagger"
	"github.com/LyricTian/gin-admin/v10/internal/wirex"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"go.uber.org/zap"
)

// The Start function initializes an injector and starts an HTTP server, with an optional pprof server,
// and returns a cleanup function.
func Start(ctx context.Context) (func(), error) {
	injector, cleanInjectorFn, err := wirex.BuildInjector(ctx)
	if err != nil {
		return nil, err
	}

	if err := injector.M.Init(ctx); err != nil {
		return nil, err
	}

	cleanHTTPServerFn, err := startHTTPServer(ctx, injector)
	if err != nil {
		return cleanInjectorFn, err
	}

	if addr := config.C.General.PprofAddr; addr != "" {
		logging.Context(ctx).Info("Pprof server is listening on " + addr)
		go func() {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				logging.Context(ctx).Error("Failed to listen pprof server", zap.Error(err))
			}
		}()
	}

	return func() {
		if cleanHTTPServerFn != nil {
			cleanHTTPServerFn()
		}

		if cleanInjectorFn != nil {
			cleanInjectorFn()
		}
	}, nil
}
