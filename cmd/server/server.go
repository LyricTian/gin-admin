package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/LyricTian/gin-admin/src"
	"github.com/LyricTian/gin-admin/src/config"
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
	swaggerDir string
)

func init() {
	flag.StringVar(&configFile, "config", "", "配置文件(.json,.yaml,.toml)")
	flag.StringVar(&configFile, "c", "", "配置文件(.json,.yaml,.toml)")
	flag.StringVar(&modelFile, "model", "", "Casbin的访问控制模型(.conf)")
	flag.StringVar(&modelFile, "m", "", "Casbin的访问控制模型(.conf)")
	flag.StringVar(&wwwDir, "www", "", "静态站点目录")
	flag.StringVar(&swaggerDir, "swagger", "", "swagger目录")
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

	casbinModelConfKey := "casbin_model_conf"
	if modelFile == "" && viper.GetString(casbinModelConfKey) == "" {
		panic("请使用-m指定casbin的访问控制模型")
	}

	if modelFile != "" {
		viper.Set(casbinModelConfKey, modelFile)
	}

	if wwwDir != "" {
		viper.Set("www", wwwDir)
	}

	if swaggerDir != "" {
		viper.Set("swagger", swaggerDir)
	}

	logger.SetVersion(VERSION)
	logger.SetTraceIDFunc(util.MustUUID)
	ctx := logger.NewTraceIDContext(context.Background(), util.MustUUID())
	span := logger.StartSpanWithCall(ctx, "主函数", "main")
	span().Printf("服务启动，运行模式：%s，版本号：%s，进程号：%d", config.GetRunMode(), VERSION, os.Getpid())

	var state int32 = 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	callback := src.Init(ctx)
	select {
	case sig := <-sc:
		atomic.StoreInt32(&state, 0)
		span().Printf("获取到退出信号[%s]", sig.String())
	}

	// 等待回调函数执行完成
	if callback != nil {
		callback()
	}

	span().Printf("服务退出")
	os.Exit(int(atomic.LoadInt32(&state)))
}
