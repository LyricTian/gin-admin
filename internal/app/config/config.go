package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/koding/multiconfig"

	"github.com/LyricTian/gin-admin/v8/pkg/util/json"
)

var (
	// C 全局配置(需要先执行MustLoad，否则拿不到配置)
	C    = new(Config)
	once sync.Once
)

// MustLoad 加载配置
func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}

		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)
	})
}

// PrintWithJSON 基于JSON格式输出配置
func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

// Config 配置参数
type Config struct {
	RunMode      string
	WWW          string
	Swagger      bool
	PrintConfig  bool
	HTTP         HTTP
	Menu         Menu
	Casbin       Casbin
	Log          Log
	LogGormHook  LogGormHook
	LogMongoHook LogMongoHook
	Root         Root
	JWTAuth      JWTAuth
	Monitor      Monitor
	Captcha      Captcha
	RateLimiter  RateLimiter
	CORS         CORS
	GZIP         GZIP
	Redis        Redis
	Gorm         Gorm
	MySQL        MySQL
	Postgres     Postgres
	Sqlite3      Sqlite3
}

// IsDebugMode 是否是debug模式
func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

// Menu 菜单配置参数
type Menu struct {
	Enable bool
	Data   string
}

// Casbin casbin配置参数
type Casbin struct {
	Enable           bool
	Debug            bool
	Model            string
	AutoLoad         bool
	AutoLoadInternal int
}

// LogHook 日志钩子
type LogHook string

// IsGorm 是否是gorm钩子
func (h LogHook) IsGorm() bool {
	return h == "gorm"
}

// Log 日志配置参数
type Log struct {
	Level         int
	Format        string
	Output        string
	OutputFile    string
	EnableHook    bool
	HookLevels    []string
	Hook          LogHook
	HookMaxThread int
	HookMaxBuffer int
}

// LogGormHook 日志gorm钩子配置
type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

// LogMongoHook 日志mongo钩子配置
type LogMongoHook struct {
	Collection string
}

// Root root用户
type Root struct {
	UserID   uint64
	UserName string
	Password string
	RealName string
}

// JWTAuth 用户认证
type JWTAuth struct {
	Enable        bool
	SigningMethod string
	SigningKey    string
	Expired       int
	Store         string
	FilePath      string
	RedisDB       int
	RedisPrefix   string
}

// HTTP http配置参数
type HTTP struct {
	Host               string
	Port               int
	CertFile           string
	KeyFile            string
	ShutdownTimeout    int
	MaxContentLength   int64
	MaxReqLoggerLength int `default:"1024"`
	MaxResLoggerLength int
}

// Monitor 监控配置参数
type Monitor struct {
	Enable    bool
	Addr      string
	ConfigDir string
}

// Captcha 图形验证码配置参数
type Captcha struct {
	Store       string
	Length      int
	Width       int
	Height      int
	RedisDB     int
	RedisPrefix string
}

// RateLimiter 请求频率限制配置参数
type RateLimiter struct {
	Enable  bool
	Count   int64
	RedisDB int
}

// CORS 跨域请求配置参数
type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

// GZIP gzip压缩
type GZIP struct {
	Enable             bool
	ExcludedExtentions []string
	ExcludedPaths      []string
}

// Redis redis配置参数
type Redis struct {
	Addr     string
	Password string
}

// Gorm gorm配置参数
type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
}

// MySQL mysql配置参数
type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

// DSN 数据库连接串
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

// Postgres postgres配置参数
type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DSN 数据库连接串
func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
}

// Sqlite3 sqlite3配置参数
type Sqlite3 struct {
	Path string
}

// DSN 数据库连接串
func (a Sqlite3) DSN() string {
	return a.Path
}
