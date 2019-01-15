package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/LyricTian/gin-admin/src"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/spf13/viper"
)

// VERSION 版本号，
// 可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=1.2.0-dev"
var VERSION = "1.2.0-dev"

var (
	configFile string
	modelFile  string
	wwwDir     string
)

func init() {
	flag.StringVar(&configFile, "config", "", "配置文件(.json,.yaml,.toml)")
	flag.StringVar(&configFile, "c", "", "配置文件(.json,.yaml,.toml)")
	flag.StringVar(&modelFile, "model", "", "Casbin的访问控制模型(.conf)")
	flag.StringVar(&modelFile, "m", "", "Casbin的访问控制模型(.conf)")
	flag.StringVar(&wwwDir, "www", "", "静态站点目录")
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

	fmt.Printf("开始运行服务，服务版本号：%s \n", VERSION)

	if modelFile != "" {
		viper.Set("casbin_model", modelFile)
	}

	if wwwDir != "" {
		viper.Set("www", wwwDir)
	}

	ctx := logger.NewTraceIDContext(context.Background(), util.MustUUID())
	callback := src.Init(ctx, VERSION)

	var state int32 = 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	select {
	case sig := <-sc:
		atomic.StoreInt32(&state, 0)
		logger.Start(ctx).Printf("获取到退出信号[%s]", sig.String())
	}

	// 等待回调函数执行完成
	if callback != nil {
		callback()
	}

	logger.Start(ctx).Printf("服务退出")
	os.Exit(int(atomic.LoadInt32(&state)))
}
