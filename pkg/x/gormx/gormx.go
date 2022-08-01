package gormx

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type ReplicasConfig struct {
	DSNs         []string
	Tables       []string
	MaxLifetime  int
	MaxIdleTime  int
	MaxOpenConns int
	MaxIdleConns int
}

type Config struct {
	Debug        bool
	DBType       string // mysql/postgres/sqlite3
	DSN          string
	MaxLifetime  int
	MaxIdleTime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
	Replicas     ReplicasConfig
}

func New(c Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(c.DBType) {
	case "mysql":
		dialector = mysql.Open(c.DSN)
	case "postgres":
		dialector = postgres.Open(c.DSN)
	default:
		_ = os.MkdirAll(filepath.Dir(c.DSN), os.ModePerm)
		dialector = sqlite.Open(c.DSN)
	}

	gconfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.TablePrefix,
			SingularTable: true,
		},
		Logger: logger.Discard,
	}

	if c.Debug {
		gconfig.Logger = logger.Default
	}

	db, err := gorm.Open(dialector, gconfig)
	if err != nil {
		return nil, err
	}

	if dsns := c.Replicas.DSNs; len(dsns) > 0 {
		dialectors := make([]gorm.Dialector, len(dsns))
		for i := range dsns {
			dialectors[i] = postgres.Open(dsns[i])
		}

		replicaTables := make([]interface{}, len(c.Replicas.Tables))
		for i, table := range c.Replicas.Tables {
			replicaTables[i] = table
		}

		err := db.Use(dbresolver.Register(
			dbresolver.Config{
				Replicas: dialectors,
			}, replicaTables...).
			SetMaxIdleConns(c.Replicas.MaxIdleConns).
			SetMaxOpenConns(c.Replicas.MaxOpenConns).
			SetConnMaxLifetime(time.Duration(c.Replicas.MaxLifetime) * time.Second).
			SetConnMaxIdleTime(time.Duration(c.Replicas.MaxIdleTime) * time.Second),
		)
		if err != nil {
			return db, err
		}
	}

	if c.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Second)

	return db, nil
}
