package mysql

import (
	"fmt"

	"github.com/LyricTian/gin-admin/src/service/mysql"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/facebookgo/inject"
	"github.com/spf13/viper"
)

// Common mysql存储模块
type Common struct {
	Demo *Demo
	User *User
	Role *Role
	Menu *Menu
}

// Init 初始化
func (a *Common) Init(g *inject.Graph, db *mysql.DB) *Common {
	a.Demo = new(Demo).Init(g, db, a)
	a.User = new(User).Init(g, db, a)
	a.Role = new(Role).Init(g, db, a)
	a.Menu = new(Menu).Init(g, db, a)
	return a
}

// TablePrefix 获取表名前缀
func (a *Common) TablePrefix() string {
	prefix := util.T(viper.GetStringMap("mysql")["table_prefix"]).String()
	if prefix != "" {
		if prefix[len(prefix)-1] != '_' {
			prefix += "_"
		}
		return prefix
	}
	return ""
}

// TableName 获取表名
func (a *Common) TableName(name string) string {
	return fmt.Sprintf("%s%s", a.TablePrefix(), name)
}
