package prometheusx

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

func (a *AdapterGin) Middleware(enable bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !enable {
			ctx.Next()
			return
		}

		b := time.Now()
		recevedBytes := 0
		if v, ok := ctx.Get("bobcatminer/req-body"); ok {
			if b, ok := v.([]byte); ok {
				recevedBytes = len(b)
			}
		}

		ctx.Next()

		if ctx.Writer.Status() == 404 {
			return
		}

		latency := float64(time.Since(b).Milliseconds())
		p := ctx.Request.URL.Path
		for _, param := range ctx.Params {
			p = strings.Replace(p, param.Value, ":"+param.Key, -1)
		}
		a.prom.Log(p, ctx.Request.Method, fmt.Sprintf("%d", ctx.Writer.Status()), float64(ctx.Writer.Size()), float64(recevedBytes), latency)
	}
}
