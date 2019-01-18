package config

import (
	"github.com/spf13/viper"
)

// GetRunMode 获取运行模式
func GetRunMode() string {
	return viper.GetString("run_mode")
}

// IsDebugMode 检查调试模式
func IsDebugMode() bool {
	return GetRunMode() == "debug"
}

// IsTestMode 检查测试模式
func IsTestMode() bool {
	return GetRunMode() == "test"
}

// IsReleaseMode 检查正式模式
func IsReleaseMode() bool {
	return GetRunMode() == "ReleaseMode"
}

// GetAuthMode 获取认证模式
func GetAuthMode() string {
	return viper.GetString("auth_mode")
}

// IsSessionAuth 会话认证
func IsSessionAuth() bool {
	return GetAuthMode() == "session"
}

// IsJWTAuth jwt认证
func IsJWTAuth() bool {
	return GetAuthMode() == "jwt"
}

// GetDBMode 获取存储模式
func GetDBMode() string {
	return viper.GetString("db_mode")
}

// IsGormDB gorm存储
func IsGormDB() bool {
	return GetDBMode() == "gorm"
}

// RootUser root用户
type RootUser struct {
	UserName string
	Password string
}

// GetRootUser 获取root用户
func GetRootUser() RootUser {
	return RootUser{
		UserName: viper.GetString("system_root_user"),
		Password: viper.GetString("system_root_password"),
	}
}

// GetDBTablePrefix 获取存储表名前缀
func GetDBTablePrefix() string {
	return viper.GetString("db_table_prefix")
}

// GetCasbinModelConf 获取casbin的模型配置文件
func GetCasbinModelConf() string {
	return viper.GetString("casbin_model_conf")
}

// GetWWWDir 获取静态站点目录
func GetWWWDir() string {
	return viper.GetString("www")
}

// GetSwaggerDir 获取swagger文档目录
func GetSwaggerDir() string {
	return viper.GetString("swagger")
}

// HTTPAddr http地址
type HTTPAddr struct {
	Host string
	Port int
}

// GetHTTPAddr 获取HTTP地址
func GetHTTPAddr() HTTPAddr {
	return HTTPAddr{
		Host: viper.GetString("http_host"),
		Port: viper.GetInt("http_port"),
	}
}
