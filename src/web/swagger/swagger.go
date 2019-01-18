/*
Package swagger 生成swagger文档

使用方式：

	go get -u -v github.com/teambition/swaggo
	cd src/web/swagger
	swaggo -d -s ./swagger.go -p ../../ -o ./swagger
*/
package swagger

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
