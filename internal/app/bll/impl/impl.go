package impl

import (
	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/bll/impl/internal"
	"go.uber.org/dig"
)

// Inject 注入bll实现
// 使用方式：
//   container := dig.New()
//   Inject(container)
//   container.Invoke(func(foo IDemo) {
//   })
func Inject(container *dig.Container) error {
	container.Provide(internal.NewTrans, dig.As(new(bll.ITrans)))
	container.Provide(internal.NewDemo, dig.As(new(bll.IDemo)))
	container.Provide(internal.NewLogin, dig.As(new(bll.ILogin)))
	container.Provide(internal.NewMenu, dig.As(new(bll.IMenu)))
	container.Provide(internal.NewRole, dig.As(new(bll.IRole)))
	container.Provide(internal.NewUser, dig.As(new(bll.IUser)))
	return nil
}
