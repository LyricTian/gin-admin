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
	C    = new(Config)
	once sync.Once
)

// Load config file (toml/json/yaml)
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

func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

type Menu struct {
	Enable bool
	Data   string
}

type Casbin struct {
	Enable           bool
	Debug            bool
	Model            string
	AutoLoad         bool
	AutoLoadInternal int
}

type LogHook string

func (h LogHook) IsGorm() bool {
	return h == "gorm"
}

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

type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

type LogMongoHook struct {
	Collection string
}

type Root struct {
	UserID   uint64
	UserName string
	Password string
	RealName string
}

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

type Monitor struct {
	Enable    bool
	Addr      string
	ConfigDir string
}

type Captcha struct {
	Store       string
	Length      int
	Width       int
	Height      int
	RedisDB     int
	RedisPrefix string
}

type RateLimiter struct {
	Enable  bool
	Count   int64
	RedisDB int
}

type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type GZIP struct {
	Enable             bool
	ExcludedExtentions []string
	ExcludedPaths      []string
}

type Redis struct {
	Addr     string
	Password string
}

type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
}

type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
}

type Sqlite3 struct {
	Path string
}

func (a Sqlite3) DSN() string {
	return a.Path
}
