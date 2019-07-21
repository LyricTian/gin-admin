package ctl

import (
	"go.uber.org/dig"
)

// Inject 注入ctl
// 使用方式：
//   container := dig.New()
//   Inject(container)
//   container.Invoke(func(foo *ctl.Demo) {
//   })
func Inject(container *dig.Container) error {
	container.Provide(NewDemo)
	container.Provide(NewLogin)
	container.Provide(NewMenu)
	container.Provide(NewRole)
	container.Provide(NewUser)
	return nil
}
