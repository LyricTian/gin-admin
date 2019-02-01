package session

import (
	"context"
	"net/http"
	"sync"
)

var (
	internalManager *Manager
	once            sync.Once
)

func manager(opt ...Option) *Manager {
	once.Do(func() {
		internalManager = NewManager(opt...)
	})
	return internalManager
}

// InitManager initialize the global session management instance
func InitManager(opt ...Option) {
	manager(opt...)
}

// Start a session and return to session storage
func Start(ctx context.Context, w http.ResponseWriter, r *http.Request) (Store, error) {
	return manager().Start(ctx, w, r)
}

// Destroy a session
func Destroy(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return manager().Destroy(ctx, w, r)
}

// Refresh a session and return to session storage
func Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) (Store, error) {
	return manager().Refresh(ctx, w, r)
}
