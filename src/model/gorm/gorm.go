package gormmodel

import (
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/model/gorm/common"
	"github.com/LyricTian/gin-admin/src/model/gorm/demo"
	"github.com/facebookgo/inject"
	"github.com/jinzhu/gorm"
)

// Init 初始化gorm存储
func Init(g *inject.Graph, db *gorm.DB) {
	g.Provide(&inject.Object{Value: model.ITrans(common.NewTrans(db)), Name: "ITrans"})
	g.Provide(&inject.Object{Value: model.IDemo(demo.NewModel(db)), Name: "IDemo"})
}
