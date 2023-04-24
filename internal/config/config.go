package config

import (
	"fmt"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/pelletier/go-toml"
)

var (
	C    = new(Config)
	once sync.Once
)

func MustLoad(name string) {
	once.Do(func() {
		tree, err := toml.LoadFile(name)
		if err != nil {
			panic(fmt.Sprintf("Failed to load config file %s: %s", name, err.Error()))
		}
		if err = tree.Unmarshal(C); err != nil {
			panic(fmt.Sprintf("Failed to unmarshal config %s: %s", name, err.Error()))
		}
		if err = C.PreLoad(); err != nil {
			panic(fmt.Sprintf("Failed to preload config %s: %s", name, err.Error()))
		}
	})
}

type Config struct {
	General    General
	Storage    Storage
	Middleware Middleware
	Util       Util
	Dictionary Dictionary
}

type General struct {
	AppName   string `default:"ginadmin"`
	DebugMode bool
	HTTP      struct {
		Addr            string `default:":8080"`
		ShutdownTimeout int    `default:"10"` // seconds
		ReadTimeout     int    `default:"60"` // seconds
		WriteTimeout    int    `default:"60"` // seconds
		IdleTimeout     int    `default:"10"` // seconds
		CertFile        string
		KeyFile         string
	}
	PprofAddr          string
	DisableSwagger     bool
	DisablePrintConfig bool
	DefaultLoginPwd    string `default:"6351623c8cef86fefabfa7da046fc619"` // abc-123
	Root               struct {
		ID       string `default:"root"`
		Username string `default:"admin"`
		Name     string `default:"Administrator"`
		Password string
	}
	LoggerConfigFile string `default:"logger.toml"`
	MenuYamlFile     string `default:"menu.yaml"`
	ConfigDir        string // From command arguments
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

type Middleware struct {
	Recovery struct {
		Skip int `default:"3"` // skip the first n stack frames
	}
	CORS struct {
		Enable                 bool
		AllowAllOrigins        bool
		AllowOrigins           []string
		AllowMethods           []string
		AllowHeaders           []string
		AllowCredentials       bool
		ExposeHeaders          []string
		MaxAge                 int
		AllowWildcard          bool
		AllowBrowserExtensions bool
		AllowWebSockets        bool
		AllowFiles             bool
	}
	Trace struct {
		SkippedPathPrefixes []string
		RequestHeaderKey    string `default:"X-Request-Id"`
		ResponseTraceKey    string `default:"X-Trace-Id"`
	}
	Logger struct {
		SkippedPathPrefixes      []string
		MaxOutputRequestBodyLen  int `default:"4096"`
		MaxOutputResponseBodyLen int `default:"1024"`
	}
	CopyBody struct {
		SkippedPathPrefixes []string
		MaxContentLen       int64 `default:"33554432"` // max content length (default 32MB)
	}
	Auth struct {
		Disable             bool
		SkippedPathPrefixes []string
		SigningMethod       string `default:"HS512"`    // HS256/HS384/HS512
		SigningKey          string `default:"XnEsT0S@"` // secret key
		OldSigningKey       string // old secret key (for migration)
		Expired             int    `default:"86400"` // seconds
		Store               struct {
			Type      string `default:"badger"` // badger/redis
			Delimiter string `default:":"`      // delimiter for key
			Badger    struct {
				Path string `default:"data/auth"`
			}
			Redis struct {
				Addr     string
				Username string
				Password string
				DB       int
			}
		}
	}
	RateLimiter struct {
		Enable              bool
		SkippedPathPrefixes []string
		Period              int // seconds
		MaxRequestsPerIP    int
		MaxRequestsPerUser  int
		Store               struct {
			Type   string // memory/redis
			Memory struct {
				Expiration      int `default:"3600"` // seconds
				CleanupInterval int `default:"60"`   // seconds
			}
			Redis struct {
				Addr     string
				Username string
				Password string
				DB       int
			}
		}
	}
	Casbin struct {
		Disable             bool
		SkippedPathPrefixes []string
		AutoLoadInterval    int    `default:"3"` // seconds
		ModelFile           string `default:"rbac_model.conf"`
	}
	Static struct {
		Dir string // Static files directory (From command arguments)
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
	return c.General.DebugMode
}

func (c *Config) String() string {
	b, err := jsoniter.MarshalIndent(c, "", "  ")
	if err != nil {
		panic("Failed to marshal config: " + err.Error())
	}
	return string(b)
}

func (c *Config) PreLoad() error {
	if addr := c.Storage.Cache.Redis.Addr; addr != "" {
		username := c.Storage.Cache.Redis.Username
		password := c.Storage.Cache.Redis.Password
		if c.Middleware.RateLimiter.Store.Redis.Addr == "" {
			c.Middleware.RateLimiter.Store.Redis.Addr = addr
			c.Middleware.RateLimiter.Store.Redis.Username = username
			c.Middleware.RateLimiter.Store.Redis.Password = password
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
