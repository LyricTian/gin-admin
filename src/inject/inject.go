package inject

import (
	"fmt"
	"time"

	mysqlModel "github.com/LyricTian/gin-admin/src/model/mysql"
	"github.com/LyricTian/gin-admin/src/service/mysql"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/casbin/casbin"
	"github.com/facebookgo/inject"
	"github.com/spf13/viper"
)

// Object 注入对象
type Object struct {
	MySQL     *mysql.DB
	Enforcer  *casbin.Enforcer
	CtlCommon *ctl.Common
}

// Init 初始化依赖注入
func Init() *Object {
	g := new(inject.Graph)

	// 注入mysql存储
	mysqlDB := initMySQL()
	new(mysqlModel.Common).Init(g, mysqlDB)

	// 注入casbin
	enforcer := casbin.NewEnforcer(viper.GetString("casbin_model"), false)
	g.Provide(&inject.Object{Value: enforcer})

	// 注入控制器
	ctlCommon := new(ctl.Common)
	g.Provide(&inject.Object{Value: ctlCommon})

	if err := g.Populate(); err != nil {
		panic("初始化依赖注入发生错误：" + err.Error())
	}

	return &Object{
		MySQL:     mysqlDB,
		Enforcer:  enforcer,
		CtlCommon: ctlCommon,
	}
}

// 初始化mysql数据库
func initMySQL() *mysql.DB {
	mysqlConfig := viper.GetStringMap("mysql")
	var opts []mysql.Option
	if v := util.T(mysqlConfig["trace"]).Bool(); v {
		opts = append(opts, mysql.SetTrace(v))
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
		mysqlConfig["username"],
		mysqlConfig["password"],
		mysqlConfig["addr"],
		mysqlConfig["database"],
	)
	opts = append(opts, mysql.SetDSN(dsn))

	if v := util.T(mysqlConfig["engine"]).String(); v != "" {
		opts = append(opts, mysql.SetEngine(v))
	}

	if v := util.T(mysqlConfig["encoding"]).String(); v != "" {
		opts = append(opts, mysql.SetEncoding(v))
	}

	if v := util.T(mysqlConfig["max_lifetime"]).Int(); v > 0 {
		opts = append(opts, mysql.SetMaxLifetime(time.Duration(v)*time.Second))
	}

	if v := util.T(mysqlConfig["max_open_conns"]).Int(); v > 0 {
		opts = append(opts, mysql.SetMaxOpenConns(v))
	}

	if v := util.T(mysqlConfig["max_idle_conns"]).Int(); v > 0 {
		opts = append(opts, mysql.SetMaxIdleConns(v))
	}

	db, err := mysql.NewDB(opts...)
	if err != nil {
		panic("初始化MySQL数据库发生错误：" + err.Error())
	}

	return db
}
