package gorm

import (
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/internal/entity"
	imodel "github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/internal/model"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
	"go.uber.org/dig"
)

// SetTablePrefix 设定表名前缀
func SetTablePrefix(prefix string) {
	entity.SetTablePrefix(prefix)
}

// AutoMigrate 自动映射数据表
func AutoMigrate(db *gormplus.DB) error {
	return db.AutoMigrate(
		new(entity.Demo),
		new(entity.User),
		new(entity.UserRole),
		new(entity.Role),
		new(entity.RoleMenu),
		new(entity.Menu),
		new(entity.MenuAction),
		new(entity.MenuResource),
	).Error
}

// Inject 注入gorm实现
// 使用方式：
//   container := dig.New()
//   Inject(container)
//   container.Invoke(func(foo IDemo) {
//   })
func Inject(container *dig.Container) error {
	container.Provide(imodel.NewTrans, dig.As(new(model.ITrans)))
	container.Provide(imodel.NewDemo, dig.As(new(model.IDemo)))
	container.Provide(imodel.NewMenu, dig.As(new(model.IMenu)))
	container.Provide(imodel.NewRole, dig.As(new(model.IRole)))
	container.Provide(imodel.NewUser, dig.As(new(model.IUser)))
	return nil
}
