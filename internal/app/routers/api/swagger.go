/*
Package api 生成swagger文档

文档规则请参考：https://github.com/swaggo/swag#declarative-comments-format

使用方式：

	go get -u github.com/swaggo/swag/cmd/swag
	swag init -g ./internal/app/routers/api/swagger.go -o ./internal/app/swagger*/
package api

// @title GinAdmin
// @version 5.0.0
// @description RBAC scaffolding based on GIN + GORM + CASBIN + Dig.
// @schemes http https
// @host 127.0.0.1:10088
// @basePath /
// @contact.name LyricTian
// @contact.email tiannianshou@gmail.com
