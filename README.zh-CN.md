[English](./README.md) | 简体中文

<h1 align="center">Gin Admin</h1>

<div align="center">
 基于 Gin + GORM + Casbin + Ant Design React 实现的RBAC权限脚手架，目的是提供一套轻量级的中后台开发框架。

[![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

![](http://store.tiannianshou.com/static/github/gin_admin_perm.png)

</div>

- :label: [在线演示地址](https://demo.tiannianshou.com) (用户名：root，密码：abc-123)（`温馨提醒：为了达到更好的演示效果，这里给出了拥有最高权限的用户，请手下留情，只操作自己新增的数据，不要动平台本身的数据！谢谢！`）
- :label: [Swagger 文档地址](https://demo.tiannianshou.com/swagger/)

## 特性

- 遵循RESTful API设计规范
- 基于Casbin的RBAC访问控制模型
- 依赖注入
- 存储分离(存储层对外采用接口的方式供业务层调用，实现了存储层的完全隔离，可以非常方便的更换存储方式)
- 支持统一的事务管理
- 日志追踪(基于[logrus](https://github.com/sirupsen/logrus)，日志钩子支持gorm)
- JWT认证(采用黑名单方式，存储支持：file/redis)
- 支持Swagger文档
- 支持跨域请求
- 支持请求频次限制
- 支持静态站点
- HTTP优雅重启
- 单元测试

## 下载并运行

### 获取代码

```
go get -v github.com/LyricTian/gin-admin/...
```

### 运行

> root用户的用户名及密码在配置文件(config/config.toml)中，默认为：root/abc-123

#### 编译并运行服务

> 也可以使用脚本运行(详情可查看`Makefile`)：make start-dev-server

```
cd $GOPATH/src/github.com/LyricTian/gin-admin
go build -o ./cmd/server/server ./cmd/server
./cmd/server/server -c ./config/config.toml -m ./config/model.conf -swagger ./src/web/swagger
```

#### 温馨提醒

1. 默认配置采用的是sqlite数据库，数据库文件在`data/gadmin.db`。如果想切换为`mysql`或`postgres`，需要先创建数据库（数据库创建脚本在`script`目录下）。
2. 日志的默认配置为写入本地文件，日志文件在`data/gadmin.log`。如果想切换到输出控制台或写入到gorm存储，可以自行切换配置。

#### 安装依赖包并运行前端

> 也可以使用脚本运行(详情可查看`Makefile`)：make start-dev-web

```
cd $GOPATH/src/github.com/LyricTian/gin-admin
cd web
npm install
npm start
```

## Swagger 文档的使用

> 文档规则请参考：[https://github.com/teambition/swaggo/wiki/Declarative-Comments-Format](https://github.com/teambition/swaggo/wiki/Declarative-Comments-Format)

### 安装工具并生成文档

```
go get -u -v github.com/teambition/swaggo
cd src/web
swaggo -s ./swagger.go -p ../ -o ./swagger
```

生成文档之后，可在浏览器中输入地址访问：[http://127.0.0.1:10088/swagger/](http://127.0.0.1:10088/swagger/)

## 项目结构

- cmd
    - server
- config
- doc
- script
- src
    - auth
    - bll
    - config
    - context
    - errors
    - inject
    - logger
    - model
        - gorm
            - entity
            - model
    - schema
    - service
    - util
    - web
        - context
        - ctl
        - middleware
        - swagger
        - test
- vendor
- web
    - config
    - public
    - src
        - assets
        - components
        - layouts
        - models
        - pages
        - services
        - utils

## MIT License

    Copyright (c) 2019 Lyric

[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/gin-admin
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/gin-admin
[godoc-url]: https://godoc.org/github.com/LyricTian/gin-admin
[godoc-image]: https://godoc.org/github.com/LyricTian/gin-admin?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
