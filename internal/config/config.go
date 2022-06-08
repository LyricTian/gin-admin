package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/LyricTian/gin-admin/v9/pkg/util/yaml"

	"github.com/koding/multiconfig"
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

func Print() {
	if C.PrintConfig {
		b, err := yaml.Marshal(C)
		if err != nil {
			fmt.Printf("[WARN] Configuration marshal yaml failed: %s \n", err.Error())
			return
		}

		fmt.Println("//-----------------   Configurations   --------------------//")
		fmt.Printf("\n %s \n", b)
		fmt.Println("//------------------------ End ---------------------------//")
	}
}

type Config struct {
	AppName     string
	RunMode     string
	ConfigDir   string
	WWW         string
	Swagger     bool
	PrintConfig bool
	HTTP        HTTP
	Casbin      Casbin
	CORS        CORS
	Log         Log
	Root        Root
	JWTAuth     JWTAuth
	Monitor     Monitor
	Captcha     Captcha
	RateLimiter RateLimiter
	Redis       Redis
	Gorm        Gorm
	MySQL       MySQL
	Postgres    Postgres
	Sqlite3     Sqlite3
	Cache       struct {
		Store   string // buntdb/redis
		Path    string // buntdb path
		RedisDB int    // redis database
	}
}

func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

func (c *Config) IsRootUser(userID string) bool {
	return c.Root.UserID == userID
}

type Casbin struct {
	Enable bool
	Debug  bool
}

type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type LogHook struct {
	Type      string   // gorm
	Levels    []string // debug/info/warn/error
	MaxBuffer int      `default:"10240"`
	MaxThread int      `default:"2"`
}

type Log struct {
	Level          int
	Format         string
	Output         string
	OutputFile     string
	RotationCount  uint
	RotationMaxAge int
	RotationSize   int
	Hooks          []LogHook
}

type Root struct {
	UserID   string
	Email    string
	Password string
}

type JWTAuth struct {
	Enable        bool
	SigningMethod string
	SigningKey    string
	OldSigningKey string
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
	MaxResLoggerLength int `default:"1024"`
}

type Monitor struct {
	Enable    bool
	Addr      string
	ConfigDir string
}

type Captcha struct {
	Enable      bool
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

type Redis struct {
	Addr     string
	Password string
}

type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxIdleTime       int
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
	Replicas struct {
		Hosts  []string
		Tables []string
	}
}

func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
}

func (a Postgres) ReplicasDSN() []string {
	list := make([]string, len(a.Replicas.Hosts))
	for i, host := range a.Replicas.Hosts {
		list[i] = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
	}
	return list
}

type Sqlite3 struct {
	Path string
}

func (a Sqlite3) DSN() string {
	return a.Path
}
