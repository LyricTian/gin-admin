# GORM store for [Session](https://github.com/go-session/session)

[![Build][Build-Status-Image]][Build-Status-Url] [![Codecov][codecov-image]][codecov-url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Quick Start

### Download and install

```bash
$ go get -u -v github.com/go-session/gorm
```

### Create file `server.go`

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	gormstore "github.com/go-session/gorm"
	"github.com/go-session/session"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	db, err := gorm.Open("sqlite3", "session.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	session.InitManager(
		session.SetStore(gormstore.NewDefaultStore(db)),
	)

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

## MIT License

    Copyright (c) 2019 Lyric

[Build-Status-Url]: https://travis-ci.org/go-session/gorm
[Build-Status-Image]: https://travis-ci.org/go-session/gorm.svg?branch=master
[codecov-url]: https://codecov.io/gh/go-session/gorm
[codecov-image]: https://codecov.io/gh/go-session/gorm/branch/master/graph/badge.svg
[reportcard-url]: https://goreportcard.com/report/github.com/go-session/gorm
[reportcard-image]: https://goreportcard.com/badge/github.com/go-session/gorm
[godoc-url]: https://godoc.org/github.com/go-session/gorm
[godoc-image]: https://godoc.org/github.com/go-session/gorm?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
