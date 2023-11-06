# [Gin](https://github.com/gin-gonic/gin)-Admin

> 基于 GIN + GORM 2.0 + Casbin 2.0 + Wire DI 的轻量级、灵活、优雅且功能齐全的 RBAC 脚手架。

[English](README.md) | 中文

[![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Features

- :scroll: 遵循 `RESTful API` 设计规范 & 基于接口的编程规范
- :house: 更加简洁的项目结构，模块化的设计，提高代码的可读性和可维护性
- :toolbox: 基于 `GIN` 框架，提供了丰富的中间件支持（JWTAuth、CORS、RequestLogger、RequestRateLimiter、TraceID、Casbin、Recover、GZIP、StaticWebsite）
- :closed_lock_with_key: 基于 `Casbin` 的 RBAC 访问控制模型
- :card_file_box: 基于 `GORM 2.0` 的数据库访问层
- :electric_plug: 基于 `WIRE` 的依赖注入 -- 依赖注入本身的作用是解决了各个模块间层级依赖繁琐的初始化过程
- :zap: 基于 `Zap & Context` 实现了日志输出，通过结合 Context 实现了统一的 TraceID/UserID 等关键字段的输出(同时支持日志钩子写入到`GORM`)
- :key: 基于 `JWT` 的用户认证
- :microscope: 基于 `Swaggo` 自动生成 `Swagger` 文档 - [Preview](http://101.42.232.163:8040/swagger/index.html)
- :test_tube: 基于 `testify` 实现了 API 的单元测试
- :100: 无状态服务，可水平扩展，提高服务的可用性 - 通过定时任务和 Redis 实现了 Casbin 的动态权限管理
- :hammer: 完善的效率工具，通过配置可以开发完整的代码模块 - [gin-admin-cli](https://github.com/gin-admin/gin-admin-cli)

![swagger](./swagger.jpeg)

## Frontend

- [基于 Ant Design React 实现的前端项目](https://github.com/gin-admin/gin-admin-frontend) - [Preview](http://101.42.232.163:8040/): admin/abc-123
- [基于 Vue.js 实现的前端项目](https://github.com/gin-admin/gin-admin-vue) - [Preview](http://101.42.232.163:8080/): admin/abc-123

## Dependencies

- [Go](https://golang.org/) 1.19+
- [Wire](github.com/google/wire) `go install github.com/google/wire/cmd/wire@latest`
- [Swag](github.com/swaggo/swag) `go install github.com/swaggo/swag/cmd/swag@latest`
- [GIN-ADMIN-CLI](https://github.com/gin-admin/gin-admin-cli) `go install github.com/gin-admin/gin-admin-cli/v10@latest`

## Quick Start

### Create a new project

```bash
gin-admin-cli new -d ~/go/src --name testapp --desc 'A test API service based on golang.' --pkg 'github.com/xxx/testapp'
```

### Start the service

```bash
cd ~/go/src/testapp
make start
# or
go run main.go start
```

### Generate a new module

> 更加详细的使用说明请参考 [gin-admin-cli](https://github.com/gin-admin/gin-admin-cli)

```bash
gin-admin-cli gen -d . -m CMS -s Article --structs-comment 'Article management'
```

### Remove a module

```bash
gin-admin-cli rm -d . -m CMS -s Article
```

### Build the service

```bash
make build
# or
go build -ldflags "-w -s -X main.VERSION=v1.0.0" -o ginadmin
```

### Generate swagger docs

```bash
make swagger
# or
swag init --parseDependency --generalInfo ./main.go --output ./internal/swagger
```

### Generate wire inject

```bash
make wire
# or
wire gen ./internal/wirex
```

## Project Layout

```text
├── cmd
│   ├── start.go
│   ├── stop.go
│   └── version.go
├── configs
│   ├── dev
│   │   ├── logging.toml
│   │   ├── middleware.toml
│   │   └── server.toml
│   ├── menu.json
│   └── rbac_model.conf
├── internal
│   ├── bootstrap
│   │   ├── bootstrap.go
│   │   ├── http.go
│   │   └── logger.go
│   ├── config
│   │   ├── config.go
│   │   ├── consts.go
│   │   ├── middleware.go
│   │   └── parse.go
│   ├── mods
│   │   ├── rbac
│   │   │   ├── api
│   │   │   ├── biz
│   │   │   ├── dal
│   │   │   ├── schema
│   │   │   ├── casbin.go
│   │   │   ├── main.go
│   │   │   └── wire.go
│   │   ├── sys
│   │   │   ├── api
│   │   │   ├── biz
│   │   │   ├── dal
│   │   │   ├── schema
│   │   │   ├── main.go
│   │   │   └── wire.go
│   │   └── mods.go
│   ├── utility
│   │   └── prom
│   │       └── prom.go
│   └── wirex
│       ├── injector.go
│       ├── wire.go
│       └── wire_gen.go
├── test
│   ├── menu_test.go
│   ├── role_test.go
│   ├── test.go
│   └── user_test.go
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
├── go.sum
└── main.go
```

## Contact Us

<div>
    <img src="https://store.zixinwangluo.cn/screenshots/gin-admin/wechat.jpeg" width="256"alt="wechat" />
    <img src="https://store.zixinwangluo.cn/screenshots/gin-admin/qqgroup.jpeg" width="256" alt="qqgroup" />
</div>

## MIT License

```text
Copyright (c) 2023 Lyric
```

[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/gin-admin
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/gin-admin
[godoc-url]: https://pkg.go.dev/github.com/LyricTian/gin-admin/v10
[godoc-image]: https://godoc.org/github.com/LyricTian/gin-admin?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
