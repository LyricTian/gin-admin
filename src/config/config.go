package config

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

var (
	lock      sync.RWMutex
	cacheData = new(sync.Map)
)

// 解析配置
func parse(key string, value interface{}) {
	if v, ok := cacheData.Load(key); ok {
		reflect.Indirect(reflect.ValueOf(value)).Set(reflect.ValueOf(v))
		return
	}

	lock.Lock()
	defer lock.Unlock()
	err := viper.UnmarshalKey(key, value)
	if err != nil {
		panic("解析配置发生错误:" + err.Error())
	}
	cacheData.Store(key, reflect.Indirect(reflect.ValueOf(value)).Interface())
}

// GetRunMode 获取运行模式
func GetRunMode() string {
	return viper.GetString("run_mode")
}

// IsDebugMode 检查调试模式
func IsDebugMode() bool {
	return GetRunMode() == "debug"
}

// IsTestMode 检查测试模式
func IsTestMode() bool {
	return GetRunMode() == "test"
}

// IsReleaseMode 检查正式模式
func IsReleaseMode() bool {
	return GetRunMode() == "release"
}

// GetCasbinModelConf 获取casbin的模型配置文件
func GetCasbinModelConf() string {
	return viper.GetString("casbin_model_conf")
}

// GetMenuJSONFile 获取存储菜单数据的JSON文件
func GetMenuJSONFile() string {
	return viper.GetString("menu_json_file")
}

// GetWWWDir 获取静态站点目录
func GetWWWDir() string {
	return viper.GetString("www")
}

// GetSwaggerDir 获取swagger文档目录
func GetSwaggerDir() string {
	return viper.GetString("swagger")
}

// GetStore 获取存储
func GetStore() string {
	return viper.GetString("store")
}

// Log 日志配置参数
type Log struct {
	Level         int    `mapstructure:"level"`
	Format        string `mapstructure:"format"`
	Output        string `mapstructure:"output"`
	OutputFile    string `mapstructure:"output_file"`
	EnableHook    bool   `mapstructure:"enable_hook"`
	Hook          string `mapstructure:"hook"`
	HookMaxThread int    `mapstructure:"hook_max_thread"`
	HookMaxBuffer int    `mapstructure:"hook_max_buffer"`
}

// GetLog 获取日志配置参数
func GetLog() Log {
	var c Log
	parse("log", &c)
	return c
}

// LogGormHook 日志gorm钩子配置
type LogGormHook struct {
	DBType       string `mapstructure:"db_type"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	Table        string `mapstructure:"table"`
}

// GetLogGormHook 获取gorm配置参数
func GetLogGormHook() LogGormHook {
	var c LogGormHook
	parse("log_gorm_hook", &c)
	return c
}

// Root root用户
type Root struct {
	UserName string `mapstructure:"user_name"`
	Password string `mapstructure:"password"`
	RealName string `mapstructure:"real_name"`
}

// GetRoot 获取root用户
func GetRoot() Root {
	var c Root
	parse("root", &c)
	return c
}

// Auth 用户认证
type Auth struct {
	SigningMethod string `mapstructure:"signing_method"`
	SigningKey    string `mapstructure:"signing_key"`
	Expired       int    `mapstructure:"expired"`
	Store         string `mapstructure:"store"`
	FilePath      string `mapstructure:"file_path"`
	RedisDB       int    `mapstructure:"redis_db"`
	RedisPrefix   string `mapstructure:"redis_prefix"`
}

// GetAuth 获取用户认证
func GetAuth() Auth {
	var c Auth
	parse("auth", &c)
	return c
}

// HTTP http配置参数
type HTTP struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

// GetHTTP 获取HTTP地址
func GetHTTP() HTTP {
	var c HTTP
	parse("http", &c)
	return c
}

// Captcha 图形验证码配置参数
type Captcha struct {
	Store       string `mapstructure:"store"`
	Length      int    `mapstructure:"length"`
	Width       int    `mapstructure:"width"`
	Height      int    `mapstructure:"height"`
	RedisDB     int    `mapstructure:"redis_db"`
	RedisPrefix string `mapstructure:"redis_prefix"`
}

// GetCaptcha 获取图形验证码配置参数
func GetCaptcha() Captcha {
	var c Captcha
	parse("captcha", &c)
	return c
}

// RateLimiter 请求频率限制配置参数
type RateLimiter struct {
	Enable  bool  `mapstructure:"enable"`
	Count   int64 `mapstructure:"count"`
	RedisDB int   `mapstructure:"redis_db"`
}

// GetRateLimiter 获取请求频率限制配置参数
func GetRateLimiter() RateLimiter {
	var c RateLimiter
	parse("rate_limiter", &c)
	return c
}

// CORS 跨域请求配置参数
type CORS struct {
	Enable           bool     `mapstructure:"enable"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// GetCORS 获取跨域请求配置参数
func GetCORS() CORS {
	var c CORS
	parse("cors", &c)
	return c
}

// Redis redis配置参数
type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
}

// GetRedis 获取redis配置参数
func GetRedis() Redis {
	var c Redis
	parse("redis", &c)
	return c
}

// Gorm gorm配置参数
type Gorm struct {
	Debug        bool   `mapstructure:"debug"`
	DBType       string `mapstructure:"db_type"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	TablePrefix  string `mapstructure:"table_prefix"`
}

// GetGorm 获取gorm配置参数
func GetGorm() Gorm {
	var c Gorm
	parse("gorm", &c)
	return c
}

// MySQL mysql配置参数
type MySQL struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	DBName     string `mapstructure:"db_name"`
	Parameters string `mapstructure:"parameters"`
}

// DSN 数据库连接串
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

// GetMySQL 获取mysql配置参数
func GetMySQL() MySQL {
	var c MySQL
	parse("mysql", &c)
	return c
}

// Postgres postgres配置参数
type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

// DSN 数据库连接串
func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password)
}

// GetPostgres 获取postgres配置参数
func GetPostgres() Postgres {
	var c Postgres
	parse("postgres", &c)
	return c
}

// Sqlite3 sqlite3配置参数
type Sqlite3 struct {
	Path string `mapstructure:"path"`
}

// DSN 数据库连接串
func (a Sqlite3) DSN() string {
	return a.Path
}

// GetSqlite3 获取sqlite3配置参数
func GetSqlite3() Sqlite3 {
	var c Sqlite3
	parse("sqlite3", &c)
	return c
}
