package ginsession

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

type (
	// ErrorHandleFunc error handling function
	ErrorHandleFunc func(*gin.Context, error)
	// Config defines the config for Session middleware
	Config struct {
		// error handling when starting the session
		ErrorHandleFunc ErrorHandleFunc
		// keys stored in the context
		StoreKey string
		// keys stored in the context
		ManageKey string
		// defines a function to skip middleware.Returning true skips processing
		// the middleware.
		Skipper func(*gin.Context) bool
	}
)

var (
	storeKey  string
	manageKey string

	// DefaultConfig is the default Recover middleware config.
	DefaultConfig = Config{
		ErrorHandleFunc: func(ctx *gin.Context, err error) {
			ctx.AbortWithError(500, err)
		},
		StoreKey:  "github.com/go-session/gin-session/store",
		ManageKey: "github.com/go-session/gin-session/manage",
		Skipper: func(_ *gin.Context) bool {
			return false
		},
	}
)

// New create a session middleware
func New(opt ...session.Option) gin.HandlerFunc {
	return NewWithConfig(DefaultConfig, opt...)
}

// NewWithConfig create a session middleware
func NewWithConfig(config Config, opt ...session.Option) gin.HandlerFunc {
	if config.ErrorHandleFunc == nil {
		config.ErrorHandleFunc = DefaultConfig.ErrorHandleFunc
	}

	manageKey = config.ManageKey
	if manageKey == "" {
		manageKey = DefaultConfig.ManageKey
	}

	storeKey = config.StoreKey
	if storeKey == "" {
		storeKey = DefaultConfig.StoreKey
	}

	manage := session.NewManager(opt...)
	return func(ctx *gin.Context) {
		if config.Skipper(ctx) {
			ctx.Next()
			return
		}

		ctx.Set(manageKey, manage)
		store, err := manage.Start(context.Background(), ctx.Writer, ctx.Request)
		if err != nil {
			config.ErrorHandleFunc(ctx, err)
			return
		}
		ctx.Set(storeKey, store)
		ctx.Next()
	}
}

// FromContext Get session storage from context
func FromContext(ctx *gin.Context) session.Store {
	return ctx.MustGet(storeKey).(session.Store)
}

// Destroy a session
func Destroy(ctx *gin.Context) error {
	return ctx.MustGet(manageKey).(*session.Manager).Destroy(context.Background(), ctx.Writer, ctx.Request)
}

// Refresh a session and return to session storage
func Refresh(ctx *gin.Context) (session.Store, error) {
	return ctx.MustGet(manageKey).(*session.Manager).Refresh(context.Background(), ctx.Writer, ctx.Request)
}
