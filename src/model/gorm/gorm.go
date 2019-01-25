package gormmodel

import (
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/model/gorm/model"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/facebookgo/inject"
)

// Init 初始化gorm存储
func Init(g *inject.Graph, db *gormplus.DB) {
	g.Provide(&inject.Object{Value: model.ITrans(gormmodel.NewTrans(db)), Name: "ITrans"})
	g.Provide(&inject.Object{Value: model.IDemo(gormmodel.InitDemo(db)), Name: "IDemo"})
	g.Provide(&inject.Object{Value: model.IMenu(gormmodel.InitMenu(db)), Name: "IMenu"})
	g.Provide(&inject.Object{Value: model.IRole(gormmodel.InitRole(db)), Name: "IRole"})
	g.Provide(&inject.Object{Value: model.IUser(gormmodel.InitUser(db)), Name: "IUser"})
}
