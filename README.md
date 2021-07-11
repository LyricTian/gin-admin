<h1 align="center">Gin Admin</h1>

<div align="center">
 基于 GIN + GORM + CASBIN + WIRE 实现的RBAC权限管理脚手架，目的是提供一套轻量的中后台开发框架，方便、快速的完成业务需求的开发。
<br/>

[![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

</div>

## 特性

- 遵循 `RESTful API` 设计规范 & 基于接口的编程规范
- 基于 `GIN` 框架，提供了丰富的中间件支持（JWTAuth、CORS、RequestLogger、RequestRateLimiter、TraceID、CasbinEnforce、Recover、GZIP）
- 基于 `Casbin` 的 RBAC 访问控制模型 -- **权限控制可以细粒度到按钮 & 接口**
- 基于 `Gorm` 的数据库存储
- 基于 `WIRE` 的依赖注入 -- 依赖注入本身的作用是解决了各个模块间层级依赖繁琐的初始化过程
- 基于 `Logrus & Context` 实现了日志输出，通过结合 Context 实现了统一的 TraceID/UserID 等关键字段的输出(同时支持日志钩子写入到`Gorm`)
- 基于 `JWT` 的用户认证 -- 基于 JWT 的黑名单验证机制
- 基于 `Swaggo` 自动生成 `Swagger` 文档 -- 独立于接口的 mock 实现
- 基于 `net/http/httptest` 标准包实现了 API 的单元测试
- 基于 `go mod` 的依赖管理(国内源可使用：<https://goproxy.cn/>)

## 依赖工具

```
go get -u github.com/cosmtrek/air
go get -u github.com/google/wire/cmd/wire
go get -u github.com/swaggo/swag/cmd/swag
```

- [air](https://github.com/cosmtrek/air) -- Live reload for Go apps
- [wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go
- [swag](https://github.com/swaggo/swag) -- Automatically generate RESTful API documentation with Swagger 2.0 for Go.

## 依赖框架

- [Gin](https://gin-gonic.com/) -- The fastest full-featured web framework for Go.
- [GORM](http://gorm.io/) -- The fantastic ORM library for Golang
- [Casbin](https://casbin.org/) -- An authorization library that supports access control models like ACL, RBAC, ABAC in Golang
- [Wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go

## 快速开始

> 也可以使用国内源：https://gitee.com/lyric/gin-admin

```bash
$ git clone https://github.com/LyricTian/gin-admin

$ cd gin-admin

# 使用AIR工具运行
$ air

# OR 基于Makefile运行
$ make start

# OR 使用go命令运行
$ go run cmd/gin-admin/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml
```

> 启动成功之后，可在浏览器中输入地址进行访问：[http://127.0.0.1:10088/swagger/index.html](http://127.0.0.1:10088/swagger/index.html)

## 生成`swagger`文档

```bash
# 基于Makefile
make swagger

# OR 使用swag命令
swag init --parseDependency --generalInfo ./cmd/${APP}/main.go --output ./internal/app/swagger
```

## 重新生成依赖注入文件

```bash
# 基于Makefile
make wire

# OR 使用wire命令
wire gen ./internal/app
```

## [gin-admin-cli](https://github.com/gin-admin/gin-admin-cli) 工具使用

### 创建项目

```bash
gin-admin-cli new -d test-gin-admin -p test-gin-admin -m
```

### 快速生成业务模块

#### 创建模板 task.yaml

```bash
name: Task
comment: 任务管理
fields:
  - name: Code
    type: string
    comment: 任务编号
    required: true
    binding_options: ""
    gorm_options: "size:50;index;"
  - name: Name
    type: string
    comment: 任务名称
    required: true
    binding_options: ""
    gorm_options: "size:50;index;"
  - name: Memo
    type: string
    comment: 任务备注
    required: false
    binding_options: ""
    gorm_options: "size:1024;"
```

#### 执行生成命令并运行

```bash
cd test-gin-admin
gin-admin-cli g -d . -p test-gin-admin -f ./task.yaml

# 生成 swagger
make swagger

# 生成依赖项
make wire

# 运行服务
make start
```

## 前端工程

- 基于[Ant Design React](https://ant.design/index-cn)版本的实现：[gin-admin-react](https://github.com/gin-admin/gin-admin-react)(也可使用国内源：[https://gitee.com/lyric/gin-admin-react](https://gitee.com/lyric/gin-admin-react))

## Questions

### OSX 下基于 `sqlite3` 运行出现：`'stdlib.h' file not found`

```bash
export CGO_ENABLED=1; export CC=gcc; go get -v -x github.com/mattn/go-sqlite3
```

## 互动交流

<div>
<img src="http://store.tiannianshou.com/screenshots/gin-admin/wechat.jpeg" width="256"alt="wechat" />
<img src="http://store.tiannianshou.com/screenshots/gin-admin/qqgroup.jpeg" width="256" alt="qqgroup" />
</div>

## MIT License

    Copyright (c) 2021 Lyric

[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/gin-admin
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/gin-admin
[godoc-url]: https://pkg.go.dev/github.com/LyricTian/gin-admin/v7
[godoc-image]: https://godoc.org/github.com/LyricTian/gin-admin?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
