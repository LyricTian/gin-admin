package bll

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
	_ = container.Provide(func() bll.ITrans { return internal.NewTrans })
	_ = container.Provide(func() bll.IDemo { return internal.NewDemo })
	_ = container.Provide(func() bll.ILogin { return internal.NewLogin })
	_ = container.Provide(func() bll.IMenu { return internal.NewMenu })
	_ = container.Provide(func() bll.IRole { return internal.NewRole })
	_ = container.Provide(func() bll.IUser { return internal.NewUser })
	return nil
}
