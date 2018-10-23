package mysql

import (
	"gin-admin/src/service/mysql"

	"github.com/facebookgo/inject"
)

// Common mysql存储模块
type Common struct {
	Demo *Demo
	Menu *Menu
}

// Init 初始化
func (a *Common) Init(g *inject.Graph, db *mysql.DB) *Common {
	a.Demo = new(Demo).Init(g, db, a)
	a.Menu = new(Menu).Init(g, db, a)
	return a
}
