package inject

import (
	"context"
	"errors"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/model/gorm"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/casbin/casbin"
	"github.com/facebookgo/inject"
)

// Object 注入对象
type Object struct {
	GormDB    *gormplus.DB
	Enforcer  *casbin.Enforcer
	CtlCommon *ctl.Common
}

// Init 初始化依赖注入
func Init(ctx context.Context) (*Object, error) {
	g := new(inject.Graph)
	obj := new(Object)

	// 注入存储层
	switch {
	case config.IsGormDB():
		db, err := getGormDB()
		if err != nil {
			return nil, err
		}
		gorm.Inject(g, db)
		obj.GormDB = db
	}

	// 注入casbin
	enforcer := casbin.NewEnforcer(config.GetCasbinModelConf(), false)
	g.Provide(&inject.Object{Value: enforcer})
	obj.Enforcer = enforcer

	// 注入控制器
	ctlCommon := new(ctl.Common)
	g.Provide(&inject.Object{Value: ctlCommon})
	obj.CtlCommon = ctlCommon

	if err := g.Populate(); err != nil {
		return nil, err
	}

	return obj, nil
}

func getGormDB() (*gormplus.DB, error) {
	cfg := config.GetGormConfig()

	var dsn string
	switch cfg.DBType {
	case "mysql":
		dsn = config.GetMySQLConfig().DSN()
	case "sqlite3":
		dsn = config.GetSqlite3Config().DSN()
	case "postgres":
		dsn = config.GetPostgresConfig().DSN()
	default:
		return nil, errors.New("unknown db")
	}

	return gormplus.New(gormplus.Config{
		Debug:        cfg.Debug,
		DBType:       cfg.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.MaxIdleConns,
		MaxLifetime:  cfg.MaxLifetime,
		MaxOpenConns: cfg.MaxOpenConns,
	})
}
