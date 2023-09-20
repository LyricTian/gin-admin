package promx

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AdapterGin struct {
	prom *PrometheusWrapper
}

func NewAdapterGin(p *PrometheusWrapper) *AdapterGin {
	return &AdapterGin{prom: p}
}

func (a *AdapterGin) Middleware(enable bool, reqKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !enable {
			ctx.Next()
			return
		}

		start := time.Now()
		recvBytes := 0
		if v, ok := ctx.Get(reqKey); ok {
			if b, ok := v.([]byte); ok {
				recvBytes = len(b)
			}
		}
		ctx.Next()
		latency := float64(time.Since(start).Milliseconds())
		p := ctx.Request.URL.Path
		for _, param := range ctx.Params {
			p = strings.Replace(p, param.Value, ":"+param.Key, -1)
		}
		a.prom.Log(p, ctx.Request.Method, fmt.Sprintf("%d", ctx.Writer.Status()), float64(ctx.Writer.Size()), float64(recvBytes), latency)
	}
}
