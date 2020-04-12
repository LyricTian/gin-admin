/*
Package app 生成swagger文档

文档规则请参考：https://github.com/swaggo/swag#declarative-comments-format

使用方式：

	go get -u github.com/swaggo/swag/cmd/swag
	swag init --generalInfo ./internal/app/swagger.go --output ./internal/app/swagger */
package app

// @title face-studio
// @version 1.0.0
// @description 提供自动化送标服务、FacePipeline参数配置、在线调参等
// @schemes http https
// @basePath /
// @contact.name NSTian
// @contact.email nstian@aibee.com
