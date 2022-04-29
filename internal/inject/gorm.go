package inject

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/dao"
	"github.com/LyricTian/gin-admin/v9/pkg/gormx"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
)

func InitGormDB(ctx context.Context) (*gorm.DB, func(), error) {
	cfg := config.C.Gorm
	db, err := NewGormDB(ctx)
	if err != nil {
		return nil, nil, err
	}

	cleanFunc := func() {
		db, err := db.DB()
		if err != nil {
			os.Stderr.WriteString("[gorm] Get db failed: " + err.Error())
			return
		}
		db.Close()
	}

	if cfg.EnableAutoMigrate {
		err = dao.AutoMigrate(db)
		if err != nil {
			return nil, cleanFunc, err
		}
	}

	return db, cleanFunc, nil
}

func NewGormDB(ctx context.Context) (*gorm.DB, error) {
	cfg := config.C
	dsn := ""
	switch cfg.Gorm.DBType {
	case "mysql":
		dsn = cfg.MySQL.DSN()
	case "sqlite3":
		dsn = cfg.Sqlite3.DSN()
		_ = os.MkdirAll(filepath.Dir(dsn), 0777)
	case "postgres":
		dsn = cfg.Postgres.DSN()
	default:
		return nil, errors.New("unknown db")
	}

	db, err := gormx.New(&gormx.Config{
		Debug:        cfg.Gorm.Debug,
		DBType:       cfg.Gorm.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.Gorm.MaxIdleConns,
		MaxLifetime:  cfg.Gorm.MaxLifetime,
		MaxIdleTime:  cfg.Gorm.MaxIdleTime,
		MaxOpenConns: cfg.Gorm.MaxOpenConns,
		TablePrefix:  cfg.Gorm.TablePrefix,
	})
	if err != nil {
		return nil, err
	}

	resolver := &dbresolver.DBResolver{}
	if dsns := cfg.Postgres.ReplicasDSN(); len(dsns) > 0 {
		dialectors := make([]gorm.Dialector, len(dsns))
		for i, dsn := range dsns {
			dialectors[i] = postgres.Open(dsn)
		}

		replicaTables := make([]interface{}, len(cfg.Postgres.Replicas.Tables))
		for i, table := range cfg.Postgres.Replicas.Tables {
			replicaTables[i] = table
		}

		resolver.Register(dbresolver.Config{
			Replicas: dialectors,
		}, replicaTables...)

		logger.WithContext(ctx).
			Infof("gorm use replicas, #tables: %v, #replicas: %s", replicaTables, dsns)
	}

	if dsns := cfg.Postgres.TileSourcesDSN(); len(dsns) > 0 {
		dialectors := make([]gorm.Dialector, len(dsns))
		for i, dsn := range dsns {
			dialectors[i] = postgres.Open(dsn)
		}

		tables := make([]interface{}, len(cfg.Postgres.TileSources.Tables))
		for i, table := range cfg.Postgres.TileSources.Tables {
			tables[i] = table
		}

		resolver.Register(dbresolver.Config{
			Sources: dialectors,
		}, tables...)

		logger.WithContext(ctx).Infof("gorm use tile sources, #tables: %v, #sources: %v", tables, dsns)
	}

	resolver.SetMaxIdleConns(cfg.Gorm.MaxIdleConns).
		SetMaxOpenConns(cfg.Gorm.MaxOpenConns).
		SetConnMaxLifetime(time.Duration(cfg.Gorm.MaxLifetime) * time.Second).
		SetConnMaxIdleTime(time.Duration(cfg.Gorm.MaxIdleTime) * time.Second)

	if err := db.Use(resolver); err != nil {
		logger.WithContext(ctx).Errorf("gorm use db resolver failed: %s", err.Error())
	}

	return db, nil
}
