<h1 align="center">Gin Admin</h1>

<div align="center">
 基于 GIN + GORM/MONGO + CASBIN + WIRE 实现的RBAC权限管理脚手架，目的是提供一套轻量的中后台开发框架，方便、快速的完成业务需求的开发。
<br/>

[![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

</div>

- [在线演示地址](http://139.129.88.71:10088) (用户名：root，密码：abc-123)（`温馨提醒：为了达到更好的演示效果，这里给出了拥有最高权限的用户，请手下留情，只操作自己新增的数据，不要动平台本身的数据！谢谢！`）
- [Swagger 文档地址](http://139.129.88.71:10088/swagger/index.html)

## 特性

- 遵循 `RESTful API` 设计规范 & 基于接口的编程规范
- 基于 `GIN` 框架，提供了丰富的中间件支持（JWTAuth、CORS、RequestLogger、RequestRateLimiter、TraceID、CasbinEnforce、Recover、GZIP）
- 基于 `Casbin` 的 RBAC 访问控制模型 -- **权限控制可以细粒度到按钮 & 接口**
- 基于 `Gorm/Mongo` 的数据库存储 -- 存储层抽象了标准的外部业务层调用接口，内部采用封闭式实现（为后续切换数据存储提供了较大的便利）
- 基于 `WIRE` 的依赖注入 -- 依赖注入本身的作用是解决了各个模块间层级依赖繁琐的初始化过程
- 基于 `Logrus & Context` 实现了日志输出，通过结合 Context 实现了统一的 TraceID/UserID 等关键字段的输出(同时支持日志钩子写入到`Gorm/Mongo`)
- 基于 `JWT` 的用户认证 -- 基于 JWT 的黑名单验证机制
- 基于 `Swaggo` 自动生成 `Swagger` 文档 -- 独立于接口的 mock 实现
- 基于 `net/http/httptest` 标准包实现了 API 的单元测试
- 基于 `go mod` 的依赖管理(国内源可使用：<https://goproxy.cn/>)
- 基于 `snowflake` 生成唯一ID

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

## 快速开始（或者使用[gin-admin-cli](https://github.com/gin-admin/gin-admin-cli)）

> 也可以使用国内源：https://gitee.com/lyric/gin-admin

```bash
$ go get -u -v github.com/LyricTian/gin-admin/v6/cmd/gin-admin
$ cd $GOPATH/src/github.com/LyricTian/gin-admin
# 使用AIR工具运行
$ air
# OR 基于Makefile运行
$ make start
# OR 使用go命令运行
$ go run cmd/gin-admin/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml
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
wire gen ./internal/app/injector
```

## 前端工程

- 基于[Ant Design React](https://ant.design/index-cn)版本的实现：[gin-admin-react](https://github.com/gin-admin/gin-admin-react)(也可使用国内源：[https://gitee.com/lyric/gin-admin-react](https://gitee.com/lyric/gin-admin-react))

## 互动交流

### 与作者对话

> 该项目是利用业余时间进行开发的，开发思路主要是源于自己的项目积累及个人思考，如果您有更好的想法和建议请与我进行沟通，一起探讨，畅聊技术人生，相互学习，一起进步。我非常期待！下面是我的微信二维码（如果此项目对您提供了帮助也可以请作者喝杯咖啡 (\*￣︶￣)，聊表心意，一起星巴克「续杯」~嘿嘿 ）：

<div>
<img src="http://store.tiannianshou.com/screenshots/gin-admin/wechat.jpeg" width="256"alt="wechat" />
<img src="http://store.tiannianshou.com/screenshots/gin-admin/we-pay.png" width="256" alt="we-pay" />
</div>

### QQ 群：1409099

<img src="http://store.tiannianshou.com/screenshots/gin-admin/qqgroup.jpeg" width="256" alt="qqgroup" />

## 付费支持

> **该框架本身的结构逻辑和实现细节点还是蛮多的，作者认为能够熟练应用此框架并快速产出作品还是需要一个摸索的阶段的。因而，作者原意贡献个人休闲娱乐时间，为大家提供服务！帮助大家快速掌握此框架，并如何使用 gin-admin-cli 快速产出功能模块，以 1\*10 的效率完成业务模块的开发。当然，在开发业务期间，遇到任何技术上的难题，作者将全心全力协助解决！再次感谢您的支持！作者会不断优化和改进此框架，更好的服务于大家！**

## MIT License

    Copyright (c) 2020 Lyric

[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/gin-admin
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/gin-admin
[godoc-url]: https://godoc.org/github.com/LyricTian/gin-admin
[godoc-image]: https://godoc.org/github.com/LyricTian/gin-admin?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
