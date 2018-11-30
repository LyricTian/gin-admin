# retry

> A tiny library for retrying failing operations with Go.

[![License][License-Image]][License-Url] [![ReportCard][ReportCard-Image]][ReportCard-Url] [![Build][Build-Status-Image]][Build-Status-Url] [![Coverage][Coverage-Image]][Coverage-Url] [![GoDoc][GoDoc-Image]][GoDoc-Url]

## Get

``` bash
go get -u github.com/LyricTian/retry
```

## Usage

``` go
package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/LyricTian/retry"
)

func main() {
	var (
		count int
		value string
	)

	err := retry.DoFunc(3, func() error {
		if count > 1 {
			value = "foo"
			return nil
		}
		count++
		return errors.New("not allowed")
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(value)
	// Output: foo
}

```

## MIT License

``` text
    Copyright (c) 2017 Lyric
```

[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[Build-Status-Url]: https://travis-ci.org/LyricTian/retry
[Build-Status-Image]: https://travis-ci.org/LyricTian/retry.svg?branch=master
[ReportCard-Url]: https://goreportcard.com/report/github.com/LyricTian/retry
[ReportCard-Image]: https://goreportcard.com/badge/github.com/LyricTian/retry
[GoDoc-Url]: https://godoc.org/github.com/LyricTian/retry
[GoDoc-Image]: https://godoc.org/github.com/LyricTian/retry?status.svg
[Coverage-Url]: https://coveralls.io/github/LyricTian/retry?branch=master
[Coverage-Image]: https://coveralls.io/repos/github/LyricTian/retry/badge.svg?branch=master
