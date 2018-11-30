package session

import (
	"context"
	"net/http"
)

// Define the key in the context
type (
	ctxResKey struct{}
	ctxReqKey struct{}
)

// returns a new Context that carries value res.
func newResContext(ctx context.Context, res http.ResponseWriter) context.Context {
	return context.WithValue(ctx, ctxResKey{}, res)
}

// FromResContext returns the ResponseWriter value stored in ctx, if any.
func FromResContext(ctx context.Context) (http.ResponseWriter, bool) {
	res, ok := ctx.Value(ctxResKey{}).(http.ResponseWriter)
	return res, ok
}

// returns a new Context that carries value req.
func newReqContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, ctxReqKey{}, req)
}

// FromReqContext returns the Request value stored in ctx, if any.
func FromReqContext(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(ctxReqKey{}).(*http.Request)
	return req, ok
}
