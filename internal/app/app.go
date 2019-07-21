package app

import (
	"context"
	"os"

	"github.com/LyricTian/gin-admin/internal/app/bll/impl"
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/casbin/casbin"
	"go.uber.org/dig"
)

type options struct {
	ConfigFile string
	ModelFile  string
	WWWDir     string
	SwaggerDir string
	Version    string
}

// Option 定义配置项
type Option func(*options)

// SetConfigFile 设定配置文件
func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

// SetModelFile 设定casbin模型配置文件
func SetModelFile(s string) Option {
	return func(o *options) {
		o.ModelFile = s
	}
}

// SetWWWDir 设定静态站点目录
func SetWWWDir(s string) Option {
	return func(o *options) {
		o.WWWDir = s
	}
}

// SetSwaggerDir 设定swagger目录
func SetSwaggerDir(s string) Option {
	return func(o *options) {
		o.SwaggerDir = s
	}
}

// SetVersion 设定版本号
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Init 应用初始化
func Init(ctx context.Context, opts ...Option) func() {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	err := config.LoadGlobalConfig(o.ConfigFile)
	handleError(err)

	cfg := config.GetGlobalConfig()

	logger.Printf(ctx, "服务启动，运行模式：%s，版本号：%s，进程号：%d", cfg.RunMode, o.Version, os.Getpid())

	if v := o.ModelFile; v != "" {
		cfg.CasbinModelConf = v
	}
	if v := o.WWWDir; v != "" {
		cfg.WWW = v
	}
	if v := o.SwaggerDir; v != "" {
		cfg.Swagger = v
	}

	loggerCall, err := InitLogger()
	handleError(err)

	err = InitMonitor()
	if err != nil {
		logger.Errorf(ctx, err.Error())
	}
	InitCaptcha()

	// 创建依赖注入容器
	container, containerCall := BuildContainer()

	// 初始化数据
	err = InitData(ctx, container)
	handleError(err)

	// 初始化HTTP服务
	httpCall := InitHTTPServer(ctx, container)
	return func() {
		if httpCall != nil {
			httpCall()
		}
		if containerCall != nil {
			containerCall()
		}
		if loggerCall != nil {
			loggerCall()
		}
	}
}

// NewEnforcer 创建casbin校验
func NewEnforcer() *casbin.Enforcer {
	cfg := config.GetGlobalConfig()
	return casbin.NewEnforcer(cfg.CasbinModelConf, false)
}

// BuildContainer 创建依赖注入容器
func BuildContainer() (*dig.Container, func()) {
	// 创建依赖注入容器
	container := dig.New()

	// 注入casbin
	container.Provide(NewEnforcer)

	// 注入认证模块
	auther, err := InitAuth()
	handleError(err)
	container.Provide(func() auth.Auther {
		return auther
	})

	// 注入存储模块
	storeCall, err := InitStore(container)
	handleError(err)

	// 注入bll
	err = impl.Inject(container)
	handleError(err)

	return container, func() {
		if auther != nil {
			auther.Release()
		}
		if storeCall != nil {
			storeCall()
		}
	}
}
