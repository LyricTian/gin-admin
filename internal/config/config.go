package config

import (
	"sync"

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
			panic("Failed to load toml file: " + err.Error())
		}

		err = tree.Unmarshal(C)
		if err != nil {
			panic("Failed to unmarshal config: " + err.Error())
		}
	})
}

type Config struct {
	General struct {
		AppName string `default:"ginadmin"`
		RunMode string `default:"debug"` // debug/test/release
		HTTP    struct {
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
		DisableInitMenu    bool
		ConfigDir          string // config directory (from command arguments)
	}
	Storage struct {
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
			MaxIdleConns int    `default:"50"`
			TablePrefix  string `default:"g_"`
			Replicas     struct {
				DSNs         []string
				Tables       []string
				MaxLifetime  int
				MaxIdleTime  int
				MaxOpenConns int
				MaxIdleConns int
			}
		}
	}
	Middleware struct {
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
		Recovery struct {
			Skip int `default:"3"` // skip the first n stack frames
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
				Type      string `default:"memory"` // memory/badger/redis
				Delimiter string `default:":"`      // delimiter for key
				Memory    struct {
					CleanupInterval int `default:"60"` // seconds
				}
				Badger struct {
					Path string
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
			Debug               bool
			SkippedPathPrefixes []string
			AutoLoadInterval    int `default:"3"` // seconds
		}
		Static struct {
			SkippedPathPrefixes []string
			Dir                 string // static directory (from command arguments)
		}
	}
	Util struct {
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
	Dictionary struct {
		RootUser struct {
			ID       string `default:"root"`
			Username string `default:"root"`
			Password string `default:"abc-123"`
			Name     string `default:"Root"`
		}
		UserCacheExpire int `default:"4"` // user cache expire in hours
	}
}

func (c *Config) IsDebug() bool {
	return c.General.RunMode == "debug"
}

func (c *Config) String() string {
	b, err := toml.Marshal(c)
	if err != nil {
		panic("Failed to marshal config with toml: " + err.Error())
	}

	return string(b)
}
