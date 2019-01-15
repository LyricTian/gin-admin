package inject

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/spf13/viper"
)

// 获取gorm存储
func getGormDB() (*gormplus.DB, error) {
	var config struct {
		Debug        bool   `mapstructure:"debug"`
		DBType       string `mapstructure:"db_type"`
		MaxLifetime  int    `mapstructure:"max_lifetime"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
	}

	err := viper.UnmarshalKey("gorm", &config)
	if err != nil {
		return nil, err
	}

	var dbConfig DBConfig
	err = viper.UnmarshalKey(config.DBType, &dbConfig)
	if err != nil {
		return nil, err
	}

	var dsn string
	switch config.DBType {
	case "mysql":
		dsn = dbConfig.MySQLDSN()
	case "sqlite3":
		dsn = dbConfig.SqliteDSN()
	case "postgres":
		dsn = dbConfig.PostgresDSN()
	default:
		return nil, errors.New("unknown db")
	}

	return gormplus.New(gormplus.Config{
		Debug:        config.Debug,
		DBType:       config.DBType,
		DSN:          dsn,
		MaxIdleConns: config.MaxIdleConns,
		MaxLifetime:  config.MaxLifetime,
		MaxOpenConns: config.MaxOpenConns,
	})
}

// DBConfig 数据库配置
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

// MySQLDSN 获取mysql连接串
func (c DBConfig) MySQLDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
	return dsn
}

// SqliteDSN sqlite连接串
func (c DBConfig) SqliteDSN() string {
	dir := filepath.Dir(c.DBName)
	if dir != "" {
		os.MkdirAll(dir, 0777)
	}

	return c.DBName
}

// PostgresDSN postgres连接串
func (c DBConfig) PostgresDSN() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		c.Host, c.Port, c.User, c.DBName, c.Password)
	return dsn
}
