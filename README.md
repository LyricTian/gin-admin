# gin-admin

> RBAC scaffolding based on GIN + Gorm 2.0 + CASBIN + WIRE (DI).

English | [中文](README_CN.md)

[![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Features

- Follow the `RESTful API` design specification
- Use `Casbin` to implement fine-grained access to the interface design
- Use `Wire` to resolve dependencies between modules
- Provides rich `Gin` middlewares (JWTAuth,CORS,RequestLogger,RequestRateLimiter,TraceID,CasbinEnforce,Recover,GZIP)
- Support `Swagger`

## Dependent Tools

```bash
go get -u github.com/cosmtrek/air
go get -u github.com/google/wire/cmd/wire
go get -u github.com/swaggo/swag/cmd/swag
```

- [air](https://github.com/cosmtrek/air) -- Live reload for Go apps
- [wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go
- [swag](https://github.com/swaggo/swag) -- Automatically generate RESTful API documentation with Swagger 2.0 for Go.

## Dependent Library

- [Gin](https://gin-gonic.com/) -- The fastest full-featured web framework for Go.
- [GORM](https://gorm.io/) -- The fantastic ORM library for Golang
- [Casbin](https://casbin.org/) -- An authorization library that supports access control models like ACL, RBAC, ABAC in Golang
- [Wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go

## Getting Started

```bash
git clone https://github.com/LyricTian/gin-admin

cd gin-admin

go run cmd/gin-admin/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml

# Or use Makefile: make start
```

> The database and table structure will be automatically created during the startup process. After the startup is successful, you can access the swagger address through the browser: [http://127.0.0.1:10088/swagger/index.html](http://127.0.0.1:10088/swagger/index.html)

### Generate `swagger` documentation

```bash
swag init --parseDependency --generalInfo ./cmd/${APP}/main.go --output ./internal/app/swagger

# Or use Makefile: make swagger
```

### Use `wire` to generate dependency injection

```bash
wire gen ./internal/app

# Or use Makefile: make wire
```

## Use the [gin-admin-cli](https://github.com/gin-admin/gin-admin-cli) tool to quickly generate modules

### Create template file: `task.yaml`

```yaml
name: Task
comment: TaskManage
fields:
  - name: Code
    type: string
    required: true
    binding_options: ""
    gorm_options: "size:50;index;"
  - name: Name
    type: string
    required: true
    binding_options: ""
    gorm_options: "size:50;index;"
  - name: Memo
    type: string
    required: false
    binding_options: ""
    gorm_options: "size:1024;"
```

### Execute `generate` command

```bash
gin-admin-cli g -d . -p github.com/LyricTian/gin-admin/v8 -f ./task.yaml

make swagger

make wire

make start
```

## Project Layout

```text
├── cmd
│   └── gin-admin
│       └── main.go       
├── configs
│   ├── config.toml       
│   ├── menu.yaml         
│   └── model.conf        
├── docs                  
├── internal
│   └── app
│       ├── api           
│       ├── config        
│       ├── contextx      
│       ├── dao           
│       ├── ginx          
│       ├── middleware    
│       ├── module        
│       ├── router        
│       ├── schema        
│       ├── service       
│       ├── swagger       
│       ├── test          
├── pkg
│   ├── auth              
│   │   └── jwtauth       
│   ├── errors            
│   ├── gormx             
│   ├── logger            
│   │   ├── hook
│   └── util              
│       ├── conv         
│       ├── hash         
│       ├── json
│       ├── snowflake
│       ├── structure
│       ├── trace
│       ├── uuid
│       └── yaml
└── scripts               
```

## Contact

<div>
<img src="http://store.tiannianshou.com/screenshots/gin-admin/wechat.jpeg" width="256"alt="wechat" />
<img src="http://store.tiannianshou.com/screenshots/gin-admin/qqgroup.jpeg" width="256" alt="qqgroup" />
</div>

## MIT License

    Copyright (c) 2021 Lyric

[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/gin-admin
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/gin-admin
[godoc-url]: https://pkg.go.dev/github.com/LyricTian/gin-admin/v8
[godoc-image]: https://godoc.org/github.com/LyricTian/gin-admin?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
