package gorm

import (
	"context"
	"os"
	"path/filepath"

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
	entity.SetTablePrefix(config.GetGorm().TablePrefix)

	db, err := NewGormDB()
	if err != nil {
		return nil, err
	}

	// 依赖注入
	g.Provide(&inject.Object{Value: model.ITrans(new(gormmodel.Trans).Init(db)), Name: "ITrans"})
	g.Provide(&inject.Object{Value: model.IDemo(new(gormmodel.Demo).Init(db)), Name: "IDemo"})
	g.Provide(&inject.Object{Value: model.IMenu(new(gormmodel.Menu).Init(db)), Name: "IMenu"})
	g.Provide(&inject.Object{Value: model.IResource(new(gormmodel.Resource).Init(db)), Name: "IResource"})
	g.Provide(&inject.Object{Value: model.IRole(new(gormmodel.Role).Init(db)), Name: "IRole"})
	g.Provide(&inject.Object{Value: model.IUser(new(gormmodel.User).Init(db)), Name: "IUser"})
	return db, nil
}

// NewGormDB 实例化gorm存储
func NewGormDB() (*gormplus.DB, error) {
	cfg := config.GetGorm()

	var dsn string
	switch cfg.DBType {
	case "mysql":
		dsn = config.GetMySQL().DSN()
	case "sqlite3":
		dsn = config.GetSqlite3().DSN()
		os.MkdirAll(filepath.Dir(dsn), 0777)
	case "postgres":
		dsn = config.GetPostgres().DSN()
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
