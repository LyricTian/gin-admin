package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/LyricTian/gin-admin/src"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

// VERSION 版本号，
// 可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=1.2.0-beta"
var VERSION = "1.2.0-beta"

var (
	configFile string
	modelFile  string
	webDir     string
)

func init() {
	flag.StringVar(&configFile, "config", "", "配置文件(.json,.yaml,.toml)")
	flag.StringVar(&configFile, "c", "", "配置文件(.json,.yaml,.toml)")
	flag.StringVar(&modelFile, "model", "", "Casbin的访问控制模型(.conf)")
	flag.StringVar(&modelFile, "m", "", "Casbin的访问控制模型(.conf)")
	flag.StringVar(&webDir, "www", "", "静态站点目录")
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

	if v := viper.GetString("casbin_model"); v == "" && modelFile == "" {
		panic("请使用-m指定Casbin的访问控制模型")
	}

	if modelFile != "" {
		viper.Set("casbin_model", modelFile)
	}

	if webDir != "" {
		viper.Set("web_dir", webDir)
	}

	ctx := util.NewTraceIDContext(context.Background(), util.MustUUID())
	httpHandler, callback := src.Init(ctx, VERSION)

	srv := &http.Server{
		Addr:           viper.GetString("http_addr"),
		Handler:        httpHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ac := make(chan error)
	// 开启HTTP监听
	go func() {
		logger.System(ctx).Infof("HTTP服务开始启动，监听地址为：[%s]", viper.GetString("http_addr"))
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			ac <- err
		}
	}()

	var state int32 = 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	select {
	case err := <-ac:
		if err != nil && atomic.LoadInt32(&state) == 1 {
			logger.System(ctx).Errorf("监听HTTP服务发生错误:%s", err.Error())
		}
	case sig := <-sc:
		atomic.StoreInt32(&state, 0)
		logger.System(ctx).Infof("获取到退出信号[%s]", sig.String())
	}

	if callback != nil {
		callback()
	}

	// 优雅关闭服务
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.System(ctx).Errorf("关闭HTTP服务发生错误:%s", err.Error())
	}

	// 退出应用
	os.Exit(int(atomic.LoadInt32(&state)))
}
