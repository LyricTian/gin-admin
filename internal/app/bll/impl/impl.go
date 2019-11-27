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
	_ = container.Provide(internal.NewTrans)
	_ = container.Provide(func(b *internal.Trans) bll.ITrans { return b })
	_ = container.Provide(internal.NewDemo)
	_ = container.Provide(func(b *internal.Demo) bll.IDemo { return b })
	_ = container.Provide(internal.NewLogin)
	_ = container.Provide(func(b *internal.Login) bll.ILogin { return b })
	_ = container.Provide(internal.NewMenu)
	_ = container.Provide(func(b *internal.Menu) bll.IMenu { return b })
	_ = container.Provide(internal.NewRole)
	_ = container.Provide(func(b *internal.Role) bll.IRole { return b })
	_ = container.Provide(internal.NewUser)
	_ = container.Provide(func(b *internal.User) bll.IUser { return b })
	return nil
}
