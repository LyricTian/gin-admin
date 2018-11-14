package main

import (
	"flag"
	"gin-admin/src"
	"gin-admin/src/logger"
	"gin-admin/src/util"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

// VERSION 服务版本号，
// 可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=1.0.1"
var VERSION = "1.0.0"

var (
	configFile string
	traceID    = util.MustUUID()
)

func init() {
	flag.StringVar(&configFile, "config", "", "配置文件(.json,.yaml,.toml)")
	flag.StringVar(&configFile, "c", "", "配置文件(.json,.yaml,.toml)")
}

func main() {
	flag.Parse()

	if configFile == "" {
		panic("请使用-c指定配置文件")
	}

	// 初始化配置文件
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件发生错误：" + err.Error())
	}

	var state int32 = 1
	ac := make(chan error)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGQUIT)

	httpHandler, closeHandle := src.Init(VERSION, traceID)

	// 开启HTTP监听
	go func() {
		logger.System(traceID).Infof("HTTP服务开始启动，端口监听在[%s]", viper.GetString("http_addr"))

		httpServer := &http.Server{
			Addr:           viper.GetString("http_addr"),
			Handler:        httpHandler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		ac <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-ac:
		if err != nil && atomic.LoadInt32(&state) == 1 {
			logger.System(traceID).Errorf("监听HTTP服务发生错误:%s", err.Error())
		}
	case sig := <-sc:
		atomic.StoreInt32(&state, 0)
		logger.System(traceID).Infof("获取到退出信号[%s]", sig.String())
	}

	if closeHandle != nil {
		closeHandle()
	}

	// 退出应用
	os.Exit(int(atomic.LoadInt32(&state)))
}
