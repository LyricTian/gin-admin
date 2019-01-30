package gorm

import (
	"context"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/model/gorm/entity"
	gormmodel "github.com/LyricTian/gin-admin/src/model/gorm/model"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/facebookgo/inject"
)

// Init 初始化gorm存储层
// 参考官方文档：http://gorm.io/zh_CN/docs/
func Init(ctx context.Context, g *inject.Graph) (*gormplus.DB, error) {
	// 设定初始值
	entity.SetTablePrefix(config.GetGormTablePrefix())

	db, err := NewGormDB()
	if err != nil {
		return nil, err
	}

	// 依赖注入
	g.Provide(&inject.Object{Value: model.ITrans(gormmodel.NewTrans(db)), Name: "ITrans"})
	g.Provide(&inject.Object{Value: model.IDemo(gormmodel.InitDemo(db)), Name: "IDemo"})
	g.Provide(&inject.Object{Value: model.IMenu(gormmodel.InitMenu(db)), Name: "IMenu"})
	g.Provide(&inject.Object{Value: model.IRole(gormmodel.InitRole(db)), Name: "IRole"})
	g.Provide(&inject.Object{Value: model.IUser(gormmodel.InitUser(db)), Name: "IUser"})
	return db, nil
}

// NewGormDB 实例化gorm存储
func NewGormDB() (*gormplus.DB, error) {
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
