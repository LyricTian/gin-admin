/*
Package ginadmin 生成swagger文档

文档规则请参考：https://github.com/teambition/swaggo/wiki/Declarative-Comments-Format

使用方式：

	go get -u -v github.com/teambition/swaggo
	cd internal/app/ginadmin
	swaggo -s ./swagger.go -p ../../../ -o ./swagger
*/
package ginadmin

import (
	// 控制器
	_ "github.com/LyricTian/gin-admin/internal/app/ginadmin/ctl"
)

// @Version 2.0.0
// @Title GinAdmin
// @Description RBAC scaffolding based on GIN + GORM + CASBIN.
// @Schemes http
// @Host 127.0.0.1:10088
// @BasePath /
// @Name LyricTian
// @Contact tiannianshou@gmail.com
// @Consumes json
// @Produces json
