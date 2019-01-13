package inject

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/LyricTian/gin-admin/src/model/gorm"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/casbin/casbin"
	"github.com/facebookgo/inject"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	// gorm存储注入
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Object 注入对象
type Object struct {
	GormDB    *gorm.DB
	Enforcer  *casbin.Enforcer
	CtlCommon *ctl.Common
}

// Init 初始化依赖注入
func Init() (*Object, error) {
	g := new(inject.Graph)
	obj := new(Object)

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

// 获取gorm存储
func getGormDB() (*gorm.DB, error) {
	var config struct {
		Debug        bool   `toml:"debug" yaml:"debug" json:"debug"`
		DBType       string `toml:"db_type" yaml:"db_type" json:"db_type"`
		MaxLifetime  int    `toml:"max_lifetime" yaml:"max_lifetime" json:"max_lifetime"`
		MaxOpenConns int    `toml:"max_open_conns" yaml:"max_open_conns" json:"max_open_conns"`
		MaxIdleConns int    `toml:"max_idle_conns" yaml:"max_idle_conns" json:"max_idle_conns"`
		TablePrefix  string `toml:"table_prefix" yaml:"table_prefix" json:"table_prefix"`
	}

	err := viper.UnmarshalKey("gorm", &config)
	if err != nil {
		return nil, err
	}

	var dsn string
	switch config.DBType {
	case "mysql":
		dsn, err = getMySQLDSN()
	case "sqlite3":
		dsn, err = getSqlite3DSN()
	case "postgres":
		dsn, err = getPostgresDSN()
	default:
		return nil, errors.New("unknown db")
	}
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(config.DBType, dsn)
	if err != nil {
		return nil, err
	}

	if config.Debug {
		db = db.Debug()
	}

	db.DB().SetMaxIdleConns(config.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	// 设定默认表名
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.TablePrefix + defaultTableName
	}

	return db, nil
}

func getMySQLDSN() (string, error) {
	var config struct {
		Host     string `toml:"host" yaml:"host" json:"host"`
		Port     int    `toml:"port" yaml:"port" json:"port"`
		User     string `toml:"user" yaml:"user" json:"user"`
		Password string `toml:"password" yaml:"password" json:"password"`
		DBName   string `toml:"db_name" yaml:"db_name" json:"db_name"`
	}

	err := viper.UnmarshalKey("mysql", &config)
	if err != nil {
		return "", err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)
	return dsn, nil
}

func getSqlite3DSN() (string, error) {
	var config struct {
		DBPath string `toml:"db_path" yaml:"db_path" json:"db_path"`
	}

	err := viper.UnmarshalKey("sqlite3", &config)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(config.DBPath)
	if dir != "" {
		os.MkdirAll(dir, 0777)
	}

	return config.DBPath, nil
}

func getPostgresDSN() (string, error) {
	var config struct {
		Host     string `toml:"host" yaml:"host" json:"host"`
		Port     int    `toml:"port" yaml:"port" json:"port"`
		User     string `toml:"user" yaml:"user" json:"user"`
		Password string `toml:"password" yaml:"password" json:"password"`
		DBName   string `toml:"db_name" yaml:"db_name" json:"db_name"`
	}

	err := viper.UnmarshalKey("postgres", &config)
	if err != nil {
		return "", err
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		config.Host, config.Port, config.User, config.DBName, config.Password)
	return dsn, nil
}
