# [Gin](https://github.com/gin-gonic/gin)-Admin

> 基于 GIN + GORM 2.0 + Casbin 2.0 + Wire DI 的轻量级、灵活、优雅且功能齐全的 RBAC 脚手架。

[English](README.md) | 中文

[![LICENSE](https://img.shields.io/github/license/LyricTian/gin-admin.svg)](https://github.com/LyricTian/gin-admin/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/LyricTian/gin-admin)](https://goreportcard.com/report/github.com/LyricTian/gin-admin)
[![GitHub release](https://img.shields.io/github/tag/LyricTian/gin-admin.svg?label=release)](https://github.com/LyricTian/gin-admin/releases)
[![GitHub release date](https://img.shields.io/github/release-date/LyricTian/gin-admin.svg)](https://github.com/LyricTian/gin-admin/releases)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/LyricTian/gin-admin)

## 功能特性

- :scroll: 遵循 `RESTful API` 设计规范 & 基于接口的编程规范
- :house: 更加简洁的项目结构，模块化的设计，提高代码的可读性和可维护性
- :rocket: 基于 `GIN` 框架，提供了丰富的中间件支持（JWTAuth、CORS、RequestLogger、RequestRateLimiter、TraceID、Casbin、Recover、GZIP、StaticWebsite）
- :closed_lock_with_key: 基于 `Casbin` 的 RBAC 访问控制模型
- :page_facing_up: 基于 `GORM 2.0` 的数据库访问层
- :electric_plug: 基于 `WIRE` 的依赖注入 -- 依赖注入本身的作用是解决了各个模块间层级依赖繁琐的初始化过程
- :memo: 基于 `Zap & Context` 实现了日志输出，通过结合 Context 实现了统一的 TraceID/UserID 等关键字段的输出(同时支持日志钩子写入到`GORM`)
- :key: 基于 `JWT` 的用户认证
- :microscope: 基于 `Swaggo` 自动生成 `Swagger` 文档 - [预览](http://101.42.232.163:8040/swagger/index.html)
- :wrench: 基于 `testify` 实现了 API 的单元测试
- :100: 无状态服务，可水平扩展，提高服务的可用性 - 通过定时任务和 Redis 实现了 Casbin 的动态权限管理
- :hammer: 完善的效率工具，通过配置可以开发完整的代码模块 - [gin-admin-cli](https://github.com/gin-admin/gin-admin-cli)

![swagger](./swagger.jpeg)

## 前端项目

- [基于 Ant Design React 实现的前端项目](https://github.com/gin-admin/gin-admin-frontend) - [预览](http://101.42.232.163:8040/): admin/abc-123
- [基于 Vue.js 实现的前端项目](https://github.com/gin-admin/gin-admin-vue) - [预览](http://101.42.232.163:8080/): admin/abc-123

## 安装依赖工具

- [Go](https://golang.org/) 1.19+
- [Wire](github.com/google/wire) `go install github.com/google/wire/cmd/wire@latest`
- [Swag](github.com/swaggo/swag) `go install github.com/swaggo/swag/cmd/swag@latest`
- [GIN-ADMIN-CLI](https://github.com/gin-admin/gin-admin-cli) `go install github.com/gin-admin/gin-admin-cli/v10@latest`

## 快速开始

### 创建新的项目

> 可通过 `gin-admin-cli help new` 查看命令的详细说明

```bash
gin-admin-cli new -d ~/go/src --name testapp --desc 'A test API service based on golang.' --pkg 'github.com/xxx/testapp' --git-url https://gitee.com/lyric/gin-admin.git
```

### 启动服务

> 通过更改 `configs/dev/server.toml` 配置文件中的 `MenuFile = "menu_cn.json"` 可以切换到中文菜单

```bash
cd ~/go/src/testapp

make start
# or
go run main.go start
```

### 编译服务

```bash
make build
# or
go build -ldflags "-w -s -X main.VERSION=v1.0.0" -o testapp
```

### 生成 Docker 镜像

```bash
sudo docker build -f ./Dockerfile -t testapp:v1.0.0 .
```

### 自动生成新的功能模块

> 支持通过配置文件生成代码，更加详细的使用说明请参考 [gin-admin-cli](https://github.com/gin-admin/gin-admin-cli)

```bash
gin-admin-cli gen -d . -m CMS --structs Article --structs-comment 'Article management'
```

### Remove a module

> 可通过 `gin-admin-cli help remove` 查看命令的详细说明

```bash
gin-admin-cli rm -d . -m CMS --structs Article
```

### 生成 Swagger 文档

> 通过 [Swag](github.com/swaggo/swag) 可以自动生成 Swagger 文档

```bash
make swagger
# or
swag init --parseDependency --generalInfo ./main.go --output ./internal/swagger
```

### 生成依赖注入代码

> 依赖注入本身的作用是解决了各个模块间层级依赖繁琐的初始化过程，通过 [Wire](github.com/google/wire) 可以自动生成依赖注入代码，简化依赖注入的过程。

```bash
make wire
# or
wire gen ./internal/wirex
```

## 项目结构概览

```text
├── cmd                             (命令行定义目录)
│   ├── start.go                    (启动命令)
│   ├── stop.go                     (停止命令)
│   └── version.go                  (版本命令)
├── configs
│   ├── dev
│   │   ├── logging.toml            (日志配置文件)
│   │   ├── middleware.toml         (中间件配置文件)
│   │   └── server.toml             (服务配置文件)
│   ├── menu.json                   (初始化菜单文件)
│   └── rbac_model.conf             (Casbin RBAC 模型配置文件)
├── internal
│   ├── bootstrap                   (初始化目录)
│   │   ├── bootstrap.go            (初始化)
│   │   ├── http.go                 (HTTP 服务)
│   │   └── logger.go               (日志服务)
│   ├── config                      (配置文件目录)
│   │   ├── config.go               (配置文件初始化)
│   │   ├── consts.go               (常量定义)
│   │   ├── middleware.go           (中间件配置)
│   │   └── parse.go                (配置文件解析)
│   ├── mods
│   │   ├── rbac                    (RBAC 模块)
│   │   │   ├── api                 (API层)
│   │   │   ├── biz                 (业务逻辑层)
│   │   │   ├── dal                 (数据访问层)
│   │   │   ├── schema              (数据模型层)
│   │   │   ├── casbin.go           (Casbin 初始化)
│   │   │   ├── main.go             (RBAC 模块入口)
│   │   │   └── wire.go             (RBAC 依赖注入初始化)
│   │   ├── sys                     (系统模块)
│   │   │   ├── api                 (API层)
│   │   │   ├── biz                 (业务逻辑层)
│   │   │   ├── dal                 (数据访问层)
│   │   │   ├── schema              (数据模型层)
│   │   │   ├── main.go             (SYS 模块入口)
│   │   │   └── wire.go             (SYS 依赖注入初始化)
│   │   └── mods.go
│   ├── utility
│   │   └── prom
│   │       └── prom.go             (Prometheus 监控，用于集成 prometheus)
│   └── wirex                       (依赖注入目录，包含了依赖组的定义和初始化)
│       ├── injector.go
│       ├── wire.go
│       └── wire_gen.go
├── test                            (单元测试目录)
│   ├── menu_test.go
│   ├── role_test.go
│   ├── test.go
│   └── user_test.go
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
├── go.sum
└── main.go                         (入口文件)
```

## 合作交流 :beers:

### 微信群

> 扫码加微信群

![wechat](https://store.zixinwangluo.cn/screenshots/gin-admin/wechat.pic.jpg-thumb320x480)

### QQ 群

![qq](https://store.zixinwangluo.cn/screenshots/gin-admin/qqgroup.jpeg-thumb320x480)
