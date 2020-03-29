package inject

// I 全局注入器
var I *Injector

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	return nil, nil, nil
}

// Injector 注入器
type Injector struct{}
