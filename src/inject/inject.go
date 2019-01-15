package inject

import (
	"github.com/LyricTian/gin-admin/src/model/gorm"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/casbin/casbin"
	"github.com/facebookgo/inject"
	"github.com/spf13/viper"
)

// Object 注入对象
type Object struct {
	GormDB    *gormplus.DB
	Enforcer  *casbin.Enforcer
	CtlCommon *ctl.Common
}

// Init 初始化依赖注入
func Init() (*Object, error) {
	g := new(inject.Graph)
	obj := new(Object)

	// 指定存储模式
	dbMode := viper.GetString("db_mode")
	switch dbMode {
	case "gorm":
		db, err := getGormDB()
		if err != nil {
			return nil, err
		}
		gormmodel.Init(g, db)
		obj.GormDB = db
	}

	// 注入casbin
	enforcer := casbin.NewEnforcer(viper.GetString("casbin_model"), false)
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
