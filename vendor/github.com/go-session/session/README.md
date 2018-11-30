# session

> A efficient, safely and easy-to-use session library for Go. 

[![Build][Build-Status-Image]][Build-Status-Url] [![Codecov][codecov-image]][codecov-url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Quick Start

### Download and install

```bash
$ go get -u -v github.com/go-session/session
```

### Create file `server.go`

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-session/session"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(context.Background(), w, r)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		store.Set("foo", "bar")
		err = store.Save()
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		http.Redirect(w, r, "/foo", 302)
	})

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(context.Background(), w, r)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		foo, ok := store.Get("foo")
		if ok {
			fmt.Fprintf(w, "foo:%s", foo)
			return
		}
		fmt.Fprint(w, "does not exist")
	})

	http.ListenAndServe(":8080", nil)
}
```

### Build and run

```bash
$ go build server.go
$ ./server
```

### Open in your web browser

<http://localhost:8080>

    foo:bar

## Features

- Easy to use
- Multi-storage support
- Multi-middleware support
- More secure, signature-based tamper-proof
- Context support

## Store Implementations

- [https://github.com/go-session/redis](https://github.com/go-session/redis) - Redis
- [https://github.com/go-session/mongo](https://github.com/go-session/mongo) - MongoDB
- [https://github.com/go-session/mysql](https://github.com/go-session/mysql) - MySQL
- [https://github.com/go-session/buntdb](https://github.com/go-session/buntdb) - [BuntDB](https://github.com/tidwall/buntdb)
- [https://github.com/go-session/cookie](https://github.com/go-session/cookie) - Cookie

## Middlewares

- [https://github.com/go-session/gin-session](https://github.com/go-session/gin-session) - [Gin](https://github.com/gin-gonic/gin)
- [https://github.com/go-session/beego-session](https://github.com/go-session/beego-session) - [Beego](https://github.com/astaxie/beego)
- [https://github.com/go-session/gear-session](https://github.com/go-session/gear-session) - [Gear](https://github.com/teambition/gear)
- [https://github.com/go-session/echo-session](https://github.com/go-session/echo-session) - [Echo](https://github.com/labstack/echo)

## MIT License

    Copyright (c) 2018 Lyric

[Build-Status-Url]: https://travis-ci.org/go-session/session
[Build-Status-Image]: https://travis-ci.org/go-session/session.svg?branch=master
[codecov-url]: https://codecov.io/gh/go-session/session
[codecov-image]: https://codecov.io/gh/go-session/session/branch/master/graph/badge.svg
[reportcard-url]: https://goreportcard.com/report/github.com/go-session/session
[reportcard-image]: https://goreportcard.com/badge/github.com/go-session/session
[godoc-url]: https://godoc.org/github.com/go-session/session
[godoc-image]: https://godoc.org/github.com/go-session/session?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
