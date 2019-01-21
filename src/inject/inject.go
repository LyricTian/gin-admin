package inject

import (
	"context"

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

	switch {
	case config.IsGormDB():
		db, err := getGormDB()
		if err != nil {
			return nil, err
		}
		gormmodel.Init(g, db)
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
