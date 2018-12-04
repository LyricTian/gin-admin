# gin-admin

> 基于GIN + Casbin + Ant Design的RBAC权限管理脚手架

## 获取

```
$ cd $GOPATH/src
$ mkdir -p github.com/LyricTian && cd github.com/LyricTian
$ git clone https://github.com/LyricTian/gin-admin.git
```

## 运行服务

```
$ cd github.com/LyricTian/gin-admin/cmd/server
$ go build -o server
$ ./server -c ../../config/config.toml
```