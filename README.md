# Gin-Admin

> 基于Gin + Casbin + Ant Design React 的RBAC权限管理脚手架

![screenshot](http://store.tiannianshou.com/static/github/gin_admin_default.png)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FLyricTian%2Fgin-admin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FLyricTian%2Fgin-admin?ref=badge_shield)

## 1. 快速开始

### 1.1 开发环境依赖

- go 1.10+
- node 8.12.0+
- yarn 1.12.3+

### 1.2 获取代码及初始化

#### 1.2.1 拉取代码

> 如果github源太慢，可以拉取国内源：<https://gitee.com/lyric/gin-admin.git>

```
git clone https://github.com/LyricTian/gin-admin.git $GOPATH/src/github.com/LyricTian/gin-admin
```

#### 1.2.2 更改数据库配置

> 编辑$GOPATH/src/github.com/LyricTian/gin-admin/config/config.toml配置文件，更改如下配置

```
# mysql数据库配置
[mysql]

# 连接地址(格式：127.0.0.1:3306)
addr = "127.0.0.1:3306"
# 用户名
username = "root"
# 密码
password = "123456"
# 数据库
database = "ginadmin"
```

#### 1.2.3 创建数据库

```
CREATE DATABASE `ginadmin` DEFAULT CHARACTER SET = `utf8mb4`;
```

#### 1.2.4 初始化数据

- 初始化菜单数据：`$GOPATH/src/github.com/LyricTian/gin-admin/script/init_menu.sql`

### 1.3 编译并运行

```
cd $GOPATH/src/github.com/LyricTian/gin-admin
```

#### 1.3.1 编译并运行服务

```
cd cmd/server
go build -o server
./server -c ../../config/config.toml -m ../../config/model.conf
```

#### 1.3.2 编译并运行前端

```
cd web
yarn
yarn start
```

#### 1.3.3 用户登录

```
在配置文件中可设定`system_root_user`配置项指定用户名及密码，默认用户名为：root，密码为：123
```

## 感谢以下框架的开源支持

- *Gin* - [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- *Casbin* - [https://github.com/casbin/casbin](https://github.com/casbin/casbin)
- *Ant Design* - [https://ant.design](https://ant.design)


## MIT License

    Copyright (c) 2018 Lyric

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FLyricTian%2Fgin-admin.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FLyricTian%2Fgin-admin?ref=badge_large)