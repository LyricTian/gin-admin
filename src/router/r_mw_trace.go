package router

import (
	"gin-admin/src/context"
	"gin-admin/src/util"
)

// TraceMiddleware 跟踪ID
func TraceMiddleware(ctx *context.Context) {
	ctx.Set(util.ContextKeyTraceID, util.UUIDString())
	ctx.Next()
}
