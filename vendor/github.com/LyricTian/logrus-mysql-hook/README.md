# An asynchronous MySQL Hook for [Logrus](https://github.com/sirupsen/logrus)

[![Build][Build-Status-Image]][Build-Status-Url] [![Codecov][codecov-image]][codecov-url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Quick Start

### Download and install

```bash
$ go get -u -v github.com/LyricTian/logrus-mysql-hook
```

### Usage

```go
import "github.com/LyricTian/logrus-mysql-hook"

// ...

mysqlHook := mysqlhook.Default(db,"log")

defer mysqlHook.Flush()

log := logrus.New()
log.AddHook(mysqlHook)
```

### Examples

```go
package main

import (
	"database/sql"
	"fmt"

	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	tableName := "t_log"
	mysqlHook := mysqlhook.Default(db, tableName)
	defer db.Exec(fmt.Sprintf("drop table %s", tableName))

	log := logrus.New()
	log.AddHook(mysqlHook)
	log.WithField("foo", "bar").Info("foo test")

	mysqlHook.Flush()

	var message string
	row := db.QueryRow(fmt.Sprintf("select message from %s", tableName))
	err = row.Scan(&message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(message)

	// Output: foo test
}
```

### Use Extra Item Examples

```go
package main

import (
	"database/sql"
	"fmt"

	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	tableName := "t_log"
	extraItems := []*mysqlhook.ExecExtraItem{
		mysqlhook.NewExecExtraItem("type", "varchar(50)"),
		mysqlhook.NewExecExtraItem("user_id", "varchar(50)"),
	}
	mysqlHook := mysqlhook.DefaultWithExtra(db, tableName, extraItems)

	defer db.Exec(fmt.Sprintf("drop table %s", tableName))

	log := logrus.New()
	log.AddHook(mysqlHook)
	log.WithField("foo", "bar").
		WithField("type", "system").
		WithField("user_id", "admin").
		Info("foo test")

	mysqlHook.Flush()

	var (
		message string
		typ     string
		userID  string
	)
	row := db.QueryRow(fmt.Sprintf("select message,type,user_id from %s", tableName))
	err = row.Scan(&message, &typ, &userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("[%s-%s]:%s\n", typ, userID, message)

	// Output: [system-admin]:foo test
}

```

## MIT License

    Copyright (c) 2018 Lyric

[Build-Status-Url]: https://travis-ci.org/LyricTian/logrus-mysql-hook
[Build-Status-Image]: https://travis-ci.org/LyricTian/logrus-mysql-hook.svg?branch=master
[codecov-url]: https://codecov.io/gh/LyricTian/logrus-mysql-hook
[codecov-image]: https://codecov.io/gh/LyricTian/logrus-mysql-hook/branch/master/graph/badge.svg
[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/logrus-mysql-hook
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/logrus-mysql-hook
[godoc-url]: https://godoc.org/github.com/LyricTian/logrus-mysql-hook
[godoc-image]: https://godoc.org/github.com/LyricTian/logrus-mysql-hook?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg