<h1 align="center">Gin Admin</h1>

<div align="center">
 基于 Gin + GORM/MONGO + Casbin + Wire 实现的RBAC权限管理脚手架，目的是提供一套轻量的中后台开发框架，方便、快速的完成业务需求的开发。
<br/>

[![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

</div>

- [在线演示地址](http://139.129.88.71:10088) (用户名：root，密码：abc-123)（`温馨提醒：为了达到更好的演示效果，这里给出了拥有最高权限的用户，请手下留情，只操作自己新增的数据，不要动平台本身的数据！谢谢！`）
- [Swagger 文档地址](http://139.129.88.71:10088/swagger/index.html)

## 特性

- 遵循 `RESTful API` 设计规范
- 基于 `GIN` 框架，提供了丰富的中间件支持（JWTAuth、CORS、RequestLogger、RequestRateLimiter、TraceID、CasbinEnforce、Recover）
- 基于 `Casbin` 的 RBAC 访问控制模型
- 基于 `GORM/Mongo` 的数据库存储 -- 存储层抽象了标准的外部业务层调用接口，内部采用封闭式实现（为后续切换数据存储提供了较大的便利）
- 基于 `WIRE` 的依赖注入 -- 依赖注入本身的作用是解决了各个模块间层级依赖繁琐的初始化过程
- 基于 `Logrus & Context` 实现了日志输出，通过结合 Context 实现了统一的 TraceID/UserID 等关键字段的输出(同时支持日志钩子写入到`Gorm/Mongo`)
- 基于 `JWT` 的用户认证 -- 基于JWT的黑名单验证机制
- 基于 `Swaggo` 自动生成 `Swagger` 文档
- 基于 `net/http/httptest` 标准包实现了 API 的单元测试

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
- [Mongo](https://github.com/mongodb/mongo-go-driver) -- The Go driver for MongoDB
- [Casbin](https://casbin.org/) -- An authorization library that supports access control models like ACL, RBAC, ABAC in Golang
- [Wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go

## 快速开始

> 也可以使用国内源：https://gitee.com/lyric/gin-admin

```
go get -v github.com/LyricTian/gin-admin/cmd/gin-admin
cd $GOPATH/src/github.com/LyricTian/gin-admin
# 使用AIR工具运行
air
# OR 基于Makefile运行
make start
# OR 使用go命令运行
go run cmd/gin-admin/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml
```

> 启动成功之后，可在浏览器中输入地址进行访问：[http://127.0.0.1:10088/swagger/index.html](http://127.0.0.1:10088/swagger/index.html)

## 生成`swagger`文档

```
# 基于Makefile
make swagger
# OR 使用swag命令
swag init --generalInfo ./internal/app/swagger.go --output ./internal/app/swagger
```

## 重新生成依赖注入文件

```
# 基于Makefile
make wire
# OR 使用wire命令
wire gen ./internal/app/initialize
```

## 项目结构概览

```
├── cmd
│   └── gin-admin # 主服务（程序入口）
├── configs # 配置文件目录(包含运行配置参数及casbin模型配置)
├── docs # 文档目录
├── internal # 内部代码模块
│   └── app # 内部应用模块入口
│       ├── api # API控制器模块
│       │   └── mock # API Mock模块(包括swagger的注释描述)
│       ├── bll # 业务逻辑模块接口
│       │   └── impl
│       │       └── bll # 业务逻辑模块接口的实现
│       ├── config # 配置参数(与config.toml一一映射)
│       ├── context # 统一上下文模块
│       ├── ginplus # gin的扩展模块
│       ├── initialize # 初始化模块（提供依赖模块的初始化函数及依赖注入的初始化）
│       ├── middleware # gin中间件模块
│       ├── model # 存储层模块接口
│       │   └── impl
│       │       └── gorm
│       │           ├── entity # 与数据库表及字段的映射实体
│       │           └── model # 存储层模块接口的gorm实现
│       │       └── mongo
│       │           ├── entity # 与数据库表及字段的映射实体
│       │           └── model # 存储层模块接口的mongo实现
│       ├── module # 内部模块间依赖的公共模块
│       ├── router # gin的路由模块
│       ├── schema # 提供Request/Response的对象模块
│       ├── swagger # swagger配置及自动生成的文件
│       └── test # API的单元测试
├── pkg # 公共模块
│   ├── auth # JWT认证模块
│   ├── errors # 统一错误处理模块
│   ├── logger # 日志模块
│   └── util # 工具库模块
├── scripts # 脚本目录
```

## 前端工程

- 基于[Ant Design React](https://ant.design/index-cn)版本的实现：[gin-admin-react](https://github.com/LyricTian/gin-admin-react)(也可使用国内源：[https://gitee.com/lyric/gin-admin-react](https://gitee.com/lyric/gin-admin-react))

## 互动交流

### 与作者对话

> 该项目是利用业余时间进行开发的，开发思路主要是源于自己的项目积累及个人思考，如果您有更好的想法和建议请与我进行沟通，一起探讨，畅聊技术人生，相互学习，一起进步。我非常期待！下面是我的微信二维码（如果此项目对您提供了帮助也可以请作者喝杯咖啡 (*￣︶￣)，聊表心意，一起星巴克「续杯」~嘿嘿 ）：

<div>
<img src="http://store.tiannianshou.com/screenshots/gin-admin/wechat.jpeg" width="256"alt="wechat" />
<img src="http://store.tiannianshou.com/screenshots/gin-admin/we-pay.png" width="256" alt="we-pay" />
</div>

### QQ 群：1409099

<img src="http://store.tiannianshou.com/screenshots/gin-admin/qqgroup.jpeg" width="256" alt="qqgroup" />

## 提供一对一服务

> 另外，作者有6年以上的工作经验，对前后端技术都有些研究。我也非常乐意通过有偿的方式为大家提供一对一专项辅导，包括：gin-admin/golang/react/vue，帮助大家快速掌握和入门！

## MIT License

    Copyright (c) 2020 Lyric

[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/gin-admin
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/gin-admin
[godoc-url]: https://godoc.org/github.com/LyricTian/gin-admin
[godoc-image]: https://godoc.org/github.com/LyricTian/gin-admin?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
