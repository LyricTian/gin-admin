package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

var (
	global *Config
)

// LoadGlobal 加载全局配置
func LoadGlobal(fpath string) error {
	c, err := Parse(fpath)
	if err != nil {
		return err
	}
	global = c
	return nil
}

// Global 获取全局配置
func Global() *Config {
	if global == nil {
		return &Config{}
	}
	return global
}

// Parse 解析配置文件
func Parse(fpath string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(fpath, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Config 配置参数
type Config struct {
	RunMode     string      `toml:"run_mode"`
	WWW         string      `toml:"www"`
	Swagger     string      `toml:"swagger"`
	Store       string      `toml:"store"`
	HTTP        HTTP        `toml:"http"`
	Menu        Menu        `toml:"menu"`
	Casbin      Casbin      `toml:"casbin"`
	Log         Log         `toml:"log"`
	LogGormHook LogGormHook `toml:"log_gorm_hook"`
	Root        Root        `toml:"root"`
	JWTAuth     JWTAuth     `toml:"jwt_auth"`
	Monitor     Monitor     `toml:"monitor"`
	Captcha     Captcha     `toml:"captcha"`
	RateLimiter RateLimiter `toml:"rate_limiter"`
	CORS        CORS        `toml:"cors"`
	Redis       Redis       `toml:"redis"`
	Gorm        Gorm        `toml:"gorm"`
	MySQL       MySQL       `toml:"mysql"`
	Postgres    Postgres    `toml:"postgres"`
	Sqlite3     Sqlite3     `toml:"sqlite3"`
}

// IsDebugMode 是否是debug模式
func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

// Menu 菜单配置参数
type Menu struct {
	Enable bool   `toml:"enable"`
	Data   string `toml:"data"`
}

// Casbin casbin配置参数
type Casbin struct {
	Enable           bool   `toml:"enable"`
	Debug            bool   `toml:"debug"`
	Model            string `toml:"model"`
	AutoLoad         bool   `toml:"auto_load"`
	AutoLoadInternal int    `toml:"auto_load_internal"`
}

// Log 日志配置参数
type Log struct {
	Level         int    `toml:"level"`
	Format        string `toml:"format"`
	Output        string `toml:"output"`
	OutputFile    string `toml:"output_file"`
	EnableHook    bool   `toml:"enable_hook"`
	Hook          string `toml:"hook"`
	HookMaxThread int    `toml:"hook_max_thread"`
	HookMaxBuffer int    `toml:"hook_max_buffer"`
}

// LogGormHook 日志gorm钩子配置
type LogGormHook struct {
	DBType       string `toml:"db_type"`
	MaxLifetime  int    `toml:"max_lifetime"`
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	Table        string `toml:"table"`
}

// Root root用户
type Root struct {
	UserName string `toml:"user_name"`
	Password string `toml:"password"`
	RealName string `toml:"real_name"`
}

// JWTAuth 用户认证
type JWTAuth struct {
	SigningMethod string `toml:"signing_method"`
	SigningKey    string `toml:"signing_key"`
	Expired       int    `toml:"expired"`
	Store         string `toml:"store"`
	FilePath      string `toml:"file_path"`
	RedisDB       int    `toml:"redis_db"`
	RedisPrefix   string `toml:"redis_prefix"`
}

// HTTP http配置参数
type HTTP struct {
	Host            string `toml:"host"`
	Port            int    `toml:"port"`
	ShutdownTimeout int    `toml:"shutdown_timeout"`
}

// Monitor 监控配置参数
type Monitor struct {
	Enable    bool   `toml:"enable"`
	Addr      string `toml:"addr"`
	ConfigDir string `toml:"config_dir"`
}

// Captcha 图形验证码配置参数
type Captcha struct {
	Store       string `toml:"store"`
	Length      int    `toml:"length"`
	Width       int    `toml:"width"`
	Height      int    `toml:"height"`
	RedisDB     int    `toml:"redis_db"`
	RedisPrefix string `toml:"redis_prefix"`
}

// RateLimiter 请求频率限制配置参数
type RateLimiter struct {
	Enable  bool  `toml:"enable"`
	Count   int64 `toml:"count"`
	RedisDB int   `toml:"redis_db"`
}

// CORS 跨域请求配置参数
type CORS struct {
	Enable           bool     `toml:"enable"`
	AllowOrigins     []string `toml:"allow_origins"`
	AllowMethods     []string `toml:"allow_methods"`
	AllowHeaders     []string `toml:"allow_headers"`
	AllowCredentials bool     `toml:"allow_credentials"`
	MaxAge           int      `toml:"max_age"`
}

// Redis redis配置参数
type Redis struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

// Gorm gorm配置参数
type Gorm struct {
	Debug             bool   `toml:"debug"`
	DBType            string `toml:"db_type"`
	MaxLifetime       int    `toml:"max_lifetime"`
	MaxOpenConns      int    `toml:"max_open_conns"`
	MaxIdleConns      int    `toml:"max_idle_conns"`
	TablePrefix       string `toml:"table_prefix"`
	EnableAutoMigrate bool   `toml:"enable_auto_migrate"`
}

// MySQL mysql配置参数
type MySQL struct {
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	DBName     string `toml:"db_name"`
	Parameters string `toml:"parameters"`
}

// DSN 数据库连接串
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

// Postgres postgres配置参数
type Postgres struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
	SSLMode  string `toml:"ssl_mode"`
}

// DSN 数据库连接串
func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
}

// Sqlite3 sqlite3配置参数
type Sqlite3 struct {
	Path string `toml:"path"`
}

// DSN 数据库连接串
func (a Sqlite3) DSN() string {
	return a.Path
}
