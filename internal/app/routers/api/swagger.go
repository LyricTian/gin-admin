/*
Package api 生成swagger文档

文档规则请参考：https://github.com/teambition/swaggo/wiki/Declarative-Comments-Format

使用方式：

	go get -u -v github.com/teambition/swaggo
	swaggo -s ./internal/app/routers/api/swagger.go -p . -o ./internal/app/swagger
*/
package api

import (
	// API控制器
	_ "github.com/LyricTian/gin-admin/internal/app/routers/api/ctl"
)

// @Version 4.0.0
// @Title GinAdmin
// @Description RBAC scaffolding based on GIN + GORM + CASBIN.
// @Schemes http,https
// @Host 127.0.0.1:10088
// @BasePath /
// @Name LyricTian
// @Contact tiannianshou@gmail.com
// @Consumes json
// @Produces json
