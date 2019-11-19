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
	_ = container.Provide(internal.NewTrans, dig.As(new(bll.ITrans)))
	_ = container.Provide(internal.NewDemo, dig.As(new(bll.IDemo)))
	_ = container.Provide(internal.NewLogin, dig.As(new(bll.ILogin)))
	_ = container.Provide(internal.NewMenu, dig.As(new(bll.IMenu)))
	_ = container.Provide(internal.NewRole, dig.As(new(bll.IRole)))
	_ = container.Provide(internal.NewUser, dig.As(new(bll.IUser)))
	return nil
}
