package gorm

import (
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/model"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/model/gorm/entity"
	gmodel "github.com/LyricTian/gin-admin/internal/app/ginadmin/model/gorm/model"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
)

// SetTablePrefix 设定表名前缀
func SetTablePrefix(prefix string) {
	entity.SetTablePrefix(prefix)
}

// AutoMigrate 自动映射数据表
func AutoMigrate(db *gormplus.DB) error {
	return db.AutoMigrate(
		new(entity.User),
		new(entity.UserRole),
		new(entity.Role),
		new(entity.RoleMenu),
		new(entity.Menu),
		new(entity.MenuAction),
		new(entity.MenuResource),
		new(entity.Demo),
	).Error
}

// NewModel 创建gorm存储，实现统一的存储接口
func NewModel(db *gormplus.DB) *model.Common {
	return &model.Common{
		Trans: gmodel.NewTrans(db),
		Demo:  gmodel.NewDemo(db),
		Menu:  gmodel.NewMenu(db),
		Role:  gmodel.NewRole(db),
		User:  gmodel.NewUser(db),
	}
}
