package api

import (
	"go.uber.org/dig"
)

// Inject 注入api
// 使用方式：
//   container := dig.New()
//   Inject(container)
//   container.Invoke(func(demo *api.Demo) {
//   })
func Inject(container *dig.Container) error {
	_ = container.Provide(NewDemo)
	_ = container.Provide(NewLogin)
	_ = container.Provide(NewMenu)
	_ = container.Provide(NewRole)
	_ = container.Provide(NewUser)
	return nil
}
