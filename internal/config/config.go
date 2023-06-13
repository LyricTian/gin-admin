package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/LyricTian/gin-admin/v10/pkg/encoding/json"
	"github.com/LyricTian/gin-admin/v10/pkg/encoding/toml"
	"github.com/LyricTian/gin-admin/v10/pkg/encoding/yaml"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/creasty/defaults"
	"github.com/pkg/errors"
)

var (
	once sync.Once
	C    = new(Config)
)

func MustLoad(dir string, names ...string) {
	once.Do(func() {
		if err := Load(dir, names...); err != nil {
			panic(err)
		}
	})
}

// The Load function loads configuration files in various formats from a directory and parses them into
// a struct.
func Load(dir string, names ...string) error {
	if err := defaults.Set(C); err != nil {
		return err
	}

	supportExts := []string{".json", ".yaml", ".yml", ".toml"}
	parseFile := func(name string) error {
		ext := filepath.Ext(name)
		if ext == "" || !strings.Contains(strings.Join(supportExts, ","), ext) {
			return nil
		}

		buf, err := os.ReadFile(name)
		if err != nil {
			return errors.Wrapf(err, "Failed to read config file %s", name)
		}

		switch ext {
		case ".json":
			err = json.Unmarshal(buf, C)
		case ".yaml", ".yml":
			err = yaml.Unmarshal(buf, C)
		case ".toml":
			err = toml.Unmarshal(buf, C)
		}
		return errors.Wrapf(err, "Failed to unmarshal config %s", name)
	}

	for _, name := range names {
		fullname := filepath.Join(dir, name)
		info, err := os.Stat(fullname)
		if err != nil {
			return errors.Wrapf(err, "Failed to get config file %s", name)
		}

		if info.IsDir() {
			err := filepath.WalkDir(fullname, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				} else if d.IsDir() {
					return nil
				}
				return parseFile(path)
			})
			if err != nil {
				return errors.Wrapf(err, "Failed to walk config dir %s", name)
			}
			continue
		}
		if err := parseFile(fullname); err != nil {
			return err
		}
	}

	return nil
}

type Config struct {
	Logger     logging.LoggerConfig
	General    General
	Storage    Storage
	Middleware Middleware
	Util       Util
	Dictionary Dictionary
}

type General struct {
	AppName            string `default:"ginadmin"`
	Debug              bool
	PprofAddr          string
	DisableSwagger     bool
	DisablePrintConfig bool
	DefaultLoginPwd    string `default:"6351623c8cef86fefabfa7da046fc619"` // abc-123
	InitMenuFile       string `default:"menu.yaml"`
	ConfigDir          string // From command arguments
	HTTP               struct {
		Addr            string `default:":8080"`
		ShutdownTimeout int    `default:"10"` // seconds
		ReadTimeout     int    `default:"60"` // seconds
		WriteTimeout    int    `default:"60"` // seconds
		IdleTimeout     int    `default:"10"` // seconds
		CertFile        string
		KeyFile         string
	}
	Root struct {
		ID       string `default:"root"`
		Username string `default:"admin"`
		Password string
		Name     string `default:"Admin"`
	}
}

type Storage struct {
	Cache struct {
		Type      string `default:"memory"` // memory/badger/redis
		Delimiter string `default:":"`      // delimiter for key
		Memory    struct {
			CleanupInterval int `default:"60"` // seconds
		}
		Badger struct {
			Path string `default:"data/cache"`
		}
		Redis struct {
			Addr     string
			Username string
			Password string
			DB       int
		}
	}
	DB struct {
		Debug        bool
		Type         string `default:"sqlite3"`                 // sqlite3/mysql/postgres
		DSN          string `default:"data/sqlite/ginadmin.db"` // database source name
		MaxLifetime  int    `default:"86400"`                   // seconds
		MaxIdleTime  int    `default:"3600"`                    // seconds
		MaxOpenConns int    `default:"100"`                     // connections
		MaxIdleConns int    `default:"50"`                      // connections
		TablePrefix  string `default:""`
		AutoMigrate  bool
		Resolver     []struct {
			DBType   string   // sqlite3/mysql/postgres
			Sources  []string // DSN
			Replicas []string // DSN
			Tables   []string
		}
	}
}

type Util struct {
	Captcha struct {
		Length    int    `default:"4"`
		Width     int    `default:"400"`
		Height    int    `default:"160"`
		CacheType string `default:"memory"` // memory/redis
		Redis     struct {
			Addr      string
			Username  string
			Password  string
			DB        int
			KeyPrefix string `default:"captcha:"`
		}
	}
}

type Dictionary struct {
	UserCacheExp int `default:"4"` // hours
}

func (c *Config) IsDebug() bool {
	return c.General.Debug
}

func (c *Config) String() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic("Failed to marshal config: " + err.Error())
	}
	return string(b)
}

func (c *Config) PreLoad() error {
	if addr := c.Storage.Cache.Redis.Addr; addr != "" {
		username := c.Storage.Cache.Redis.Username
		password := c.Storage.Cache.Redis.Password
		if c.Middleware.RateLimiter.Store.Type == "redis" &&
			c.Middleware.RateLimiter.Store.Redis.Addr == "" {
			c.Middleware.RateLimiter.Store.Redis.Addr = addr
			c.Middleware.RateLimiter.Store.Redis.Username = username
			c.Middleware.RateLimiter.Store.Redis.Password = password
		}
		if c.Middleware.Auth.Store.Type == "redis" &&
			c.Middleware.Auth.Store.Redis.Addr == "" {
			c.Middleware.Auth.Store.Redis.Addr = addr
			c.Middleware.Auth.Store.Redis.Username = username
			c.Middleware.Auth.Store.Redis.Password = password
		}
	}
	return nil
}

func (c *Config) Print() {
	if c.General.DisablePrintConfig {
		return
	}
	fmt.Println("// ----------------------- Load configurations start ------------------------")
	fmt.Println(c.String())
	fmt.Println("// ----------------------- Load configurations end --------------------------")
}
