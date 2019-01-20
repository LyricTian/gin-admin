/*
Package web 生成swagger文档

文档规则请参考：https://github.com/teambition/swaggo/wiki/Declarative-Comments-Format

使用方式：

	go get -u -v github.com/teambition/swaggo
	cd src/web
	swaggo -s ./swagger.go -p ../ -o ./swagger
*/
package web

import (
	// 控制器
	_ "github.com/LyricTian/gin-admin/src/web/ctl"
)

// @Version 1.2.0-dev
// @Title GinAdmin
// @Description RBAC scaffolding based on GIN + GORM + CASBIN + Ant Design React.
// @Schemes http
// @Host 127.0.0.1:8086
// @BasePath /
// @Name LyricTian
// @Contact tiannianshou@gmail.com
// @Consumes json
// @Produces json
