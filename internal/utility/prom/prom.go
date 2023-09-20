package prom

import (
	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/pkg/promx"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"github.com/gin-gonic/gin"
)

var (
	Ins           *promx.PrometheusWrapper
	GinMiddleware gin.HandlerFunc
)

func Init() {
	logMethod := make(map[string]struct{})
	logAPI := make(map[string]struct{})
	for _, m := range config.C.Util.Prometheus.LogMethods {
		logMethod[m] = struct{}{}
	}
	for _, a := range config.C.Util.Prometheus.LogApis {
		logAPI[a] = struct{}{}
	}
	Ins = promx.NewPrometheusWrapper(&promx.Config{
		Enable:         config.C.Util.Prometheus.Enable,
		App:            config.C.General.AppName,
		ListenPort:     config.C.Util.Prometheus.Port,
		BasicUserName:  config.C.Util.Prometheus.BasicUsername,
		BasicPassword:  config.C.Util.Prometheus.BasicPassword,
		LogApi:         logAPI,
		LogMethod:      logMethod,
		Objectives:     map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		DefaultCollect: config.C.Util.Prometheus.DefaultCollect,
	})
	GinMiddleware = promx.NewAdapterGin(Ins).Middleware(config.C.Util.Prometheus.Enable, util.ReqBodyKey)
}
