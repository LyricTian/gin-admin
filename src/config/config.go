package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// 解析配置
func parse(key string, value interface{}) {
	err := viper.UnmarshalKey(key, value)
	if err != nil {
		panic("解析配置发生错误:" + err.Error())
	}
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

// GetAuthMode 获取认证模式
func GetAuthMode() string {
	return viper.GetString("auth_mode")
}

// IsSessionAuth 会话认证
func IsSessionAuth() bool {
	return GetAuthMode() == "session"
}

// IsJWTAuth jwt认证
func IsJWTAuth() bool {
	return GetAuthMode() == "jwt"
}

// GetDBMode 获取存储模式
func GetDBMode() string {
	return viper.GetString("db_mode")
}

// IsGormDB gorm存储
func IsGormDB() bool {
	return GetDBMode() == "gorm"
}

// GetCasbinModelConf 获取casbin的模型配置文件
func GetCasbinModelConf() string {
	return viper.GetString("casbin_model_conf")
}

// GetWWWDir 获取静态站点目录
func GetWWWDir() string {
	return viper.GetString("www")
}

// GetSwaggerDir 获取swagger文档目录
func GetSwaggerDir() string {
	return viper.GetString("swagger")
}

// RootUser root用户
type RootUser struct {
	UserName string `mapstructure:"user_name"`
	Password string `mapstructure:"password"`
	RealName string `mapstructure:"real_name"`
}

// GetRootUser 获取root用户
func GetRootUser() RootUser {
	var config RootUser
	parse("root_user", &config)
	return config
}

// HTTPConfig http配置参数
type HTTPConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

// GetHTTPConfig 获取HTTP地址
func GetHTTPConfig() HTTPConfig {
	var config HTTPConfig
	parse("http", &config)
	return config
}

// CaptchaConfig 图形验证码配置参数
type CaptchaConfig struct {
	Store  string `mapstructure:"store"`
	Length int    `mapstructure:"length"`
	Width  int    `mapstructure:"width"`
	Height int    `mapstructure:"height"`
}

// GetCaptchaConfig 获取图形验证码配置参数
func GetCaptchaConfig() CaptchaConfig {
	var config CaptchaConfig
	parse("captcha", &config)
	return config
}

// IsCaptchaRedisStore 图形验证码存储是否是redis存储
func IsCaptchaRedisStore() bool {
	return GetCaptchaConfig().Store == "redis"
}

// RedisConfig redis配置参数
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

// GetRedisConfig 获取redis配置参数
func GetRedisConfig() RedisConfig {
	var config RedisConfig
	parse("redis", &config)
	return config
}

// LogConfig 日志配置参数
type LogConfig struct {
	Level         int    `mapstructure:"level"`
	Format        string `mapstructure:"format"`
	EnableHook    bool   `mapstructure:"enable_hook"`
	HookMaxThread int    `mapstructure:"hook_max_thread"`
	HookMaxBuffer int    `mapstructure:"hook_max_buffer"`
}

// GetLogConfig 获取日志配置参数
func GetLogConfig() LogConfig {
	var config LogConfig
	parse("log", &config)
	return config
}

// GormConfig gorm配置参数
type GormConfig struct {
	Debug        bool   `mapstructure:"debug"`
	DBType       string `mapstructure:"db_type"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	TablePrefix  string `mapstructure:"table_prefix"`
}

// GetGormConfig 获取gorm配置参数
func GetGormConfig() GormConfig {
	var config GormConfig
	parse("gorm", &config)
	return config
}

// GetGormTablePrefix 获取gorm表名前缀
func GetGormTablePrefix() string {
	return GetGormConfig().TablePrefix
}

// SessionConfig 会话配置参数
type SessionConfig struct {
	HeaderName  string `mapstructure:"header_name"`
	Sign        string `mapstructure:"sign"`
	Expired     int64  `mapstructure:"expired"`
	EnableStore bool   `mapstructure:"enable_store"`
}

// GetSessionConfig 获取会话配置参数
func GetSessionConfig() SessionConfig {
	var config SessionConfig
	parse("session", &config)
	return config
}

// MySQLConfig mysql配置参数
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

// DSN 数据库连接串
func (a MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		a.User, a.Password, a.Host, a.Port, a.DBName)
}

// GetMySQLConfig 获取mysql配置参数
func GetMySQLConfig() MySQLConfig {
	var config MySQLConfig
	parse("mysql", &config)
	return config
}

// PostgresConfig postgres配置参数
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

// DSN 数据库连接串
func (a PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password)
}

// GetPostgresConfig 获取postgres配置参数
func GetPostgresConfig() PostgresConfig {
	var config PostgresConfig
	parse("postgres", &config)
	return config
}

// Sqlite3Config sqlite3配置参数
type Sqlite3Config struct {
	Path string `mapstructure:"path"`
}

// DSN 数据库连接串
func (a Sqlite3Config) DSN() string {
	return a.Path
}

// GetSqlite3Config 获取sqlite3配置参数
func GetSqlite3Config() Sqlite3Config {
	var config Sqlite3Config
	parse("sqlite3", &config)
	return config
}
