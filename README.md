# [Gin](https://github.com/gin-gonic/gin)-Admin

> 基于 Golang + Gin + GORM 2.0 + Casbin 2.0 + Wire DI 的轻量级、灵活、优雅且功能齐全的 RBAC 脚手架。

[English](README_EN.md) | 中文

[![LICENSE](https://img.shields.io/github/license/LyricTian/gin-admin.svg)](https://github.com/LyricTian/gin-admin/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/LyricTian/gin-admin)](https://goreportcard.com/report/github.com/LyricTian/gin-admin)
[![GitHub release](https://img.shields.io/github/tag/LyricTian/gin-admin.svg?label=release)](https://github.com/LyricTian/gin-admin/releases)
[![GitHub release date](https://img.shields.io/github/release-date/LyricTian/gin-admin.svg)](https://github.com/LyricTian/gin-admin/releases)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/LyricTian/gin-admin)

## 功能特性

- :scroll: 优雅实现 `RESTful API`，采用接口化编程范式，让您的 API 设计更加专业规范
- :house: 采用清晰简洁的模块化架构，让代码结构一目了然，维护升级更轻松自如
- :rocket: 基于高性能 `GIN` 框架，集成丰富实用的中间件（身份认证、跨域、日志、限流、链路追踪、权限控制、容错、压缩等），助您快速构建企业级应用
- :closed_lock_with_key: 集成业界领先的 `Casbin` 权限框架，灵活精准的 RBAC 权限控制让安全防护固若金汤
- :page_facing_up: 基于功能强大的 `GORM 2.0` ORM 框架，优雅处理数据库操作，大幅提升开发效率
- :electric_plug: 创新采用 `WIRE` 依赖注入，革命性地简化模块依赖关系，让代码更加优雅解耦
- :memo: 基于高性能 `Zap` 日志框架，配合 Context 链路追踪，让系统运行状态清晰透明，问题排查无所遁形
- :key: 整合久经考验的 `JWT` 认证机制，让用户身份验证更加安全可靠
- :microscope: 自动集成 `Swagger` 接口文档，API 文档实时更新，开发调试更轻松 - [在线体验](https://demo.ginadmin.top/swagger/index.html)
- :wrench: 完善的单元测试体系，基于 `testify` 框架保障系统质量，让 bug 无处藏身
- :100: 采用无状态设计，支持水平扩展，搭配 Redis 实现动态权限管理，让您的系统轻松应对高并发
- :hammer: 开发者福音！配套强大的脚手架工具 [gin-admin-cli](https://github.com/gin-admin/gin-admin-cli)，让您的开发工作事半功倍

![demo](./demo.png)
![swagger](./swagger.png)

## 前端项目

- [基于 Ant Design React 实现的前端项目](https://github.com/gin-admin/gin-admin-frontend)
- [基于 Vue.js 实现的前端项目](https://github.com/gin-admin/gin-admin-vue)

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
docker build -f ./Dockerfile -t testapp:v1.0.0 .
```

### 代码生成

> 可通过 `gin-admin-cli help gen` 查看命令的详细说明

#### 准备配置文件 `dictionary.yaml`

```yaml
- name: Dictionary
  comment: 字典管理
  disable_pagination: true
  fill_gorm_commit: true
  fill_router_prefix: true
  tpl_type: "tree"
  fields:
    - name: Code
      type: string
      comment: Code of dictionary (unique for same parent)
      gorm_tag: "size:32;"
      form:
        binding_tag: "required,max=32"
    - name: Name
      type: string
      comment: Display name of dictionary
      gorm_tag: "size:128;index"
      query:
        name: LikeName
        in_query: true
        form_tag: name
        op: LIKE
      form:
        binding_tag: "required,max=128"
    - name: Description
      type: string
      comment: Details about dictionary
      gorm_tag: "size:1024"
      form: {}
    - name: Sequence
      type: int
      comment: Sequence for sorting
      gorm_tag: "index;"
      order: DESC
      form: {}
    - name: Status
      type: string
      comment: Status of dictionary (disabled, enabled)
      gorm_tag: "size:20;index"
      query: {}
      form:
        binding_tag: "required,oneof=disabled enabled"
```

```bash
gin-admin-cli gen -d . -m SYS -c dictionary.yaml
```

### 删除模块

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
│   ├── start.go                    (启动命令)
│   ├── stop.go                     (停止命令)
│   └── version.go                  (版本命令)
├── configs
│   ├── dev
│   │   ├── logging.toml            (日志配置文件)
│   │   ├── middleware.toml         (中间件配置文件)
│   │   └── server.toml             (服务配置文件)
│   ├── menu.json                   (初始化菜单文件)
│   └── rbac_model.conf             (Casbin RBAC 模型配置文件)
├── internal
│   ├── bootstrap                   (初始化目录)
│   │   ├── bootstrap.go            (初始化)
│   │   ├── http.go                 (HTTP 服务)
│   │   └── logger.go               (日志服务)
│   ├── config                      (配置文件目录)
│   │   ├── config.go               (配置文件初始化)
│   │   ├── consts.go               (常量定义)
│   │   ├── middleware.go           (中间件配置)
│   │   └── parse.go                (配置文件解析)
│   ├── mods
│   │   ├── rbac                    (RBAC 模块)
│   │   │   ├── api                 (API层)
│   │   │   ├── biz                 (业务逻辑层)
│   │   │   ├── dal                 (数据访问层)
│   │   │   ├── schema              (数据模型层)
│   │   │   ├── casbin.go           (Casbin 初始化)
│   │   │   ├── main.go             (RBAC 模块入口)
│   │   │   └── wire.go             (RBAC 依赖注入初始化)
│   │   └── mods.go
│   ├── utility
│   │   └── prom
│   │       └── prom.go             (Prometheus 监控，用于集成 prometheus)
│   └── wirex                       (依赖注入目录，包含了依赖组的定义和初始化)
│       ├── injector.go
│       ├── wire.go
│       └── wire_gen.go
├── pkg                             (公共包目录)
│   ├── cachex                      (缓存包)
│   ├── crypto                      (加密包)
│   │   ├── aes                     (AES加密)
│   │   ├── hash                    (哈希加密)
│   │   └── rand                    (随机数)
│   ├── encoding                    (编码包)
│   │   ├── json                    (JSON编码)
│   │   ├── toml                    (TOML编码)
│   │   └── yaml                    (YAML编码)
│   ├── errors                      (错误处理包)
│   ├── gormx                       (Gorm扩展包)
│   ├── jwtx                        (JWT包)
│   ├── logging                     (日志包)
│   ├── mail                        (邮件包)
│   ├── middleware                  (中间件包)
│   ├── oss                         (对象存储包)
│   ├── promx                       (Prometheus包)
│   └── util                        (工具包)
├── test                            (单元测试目录)
│   ├── menu_test.go
│   ├── role_test.go
│   ├── test.go
│   └── user_test.go
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
├── go.sum
└── main.go                         (入口文件)
```

## 合作交流 :beers:

### 微信

> 扫码加微信群

![wechat](./wechat.jpg)

### QQ

![qq](./qqgroup.jpg)

## License

Copyright (c) 2024 Lyric

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
