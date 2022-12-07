# zgin
[![Go Report Card](https://goreportcard.com/badge/github.com/ZYallers/zgin)](https://goreportcard.com/report/github.com/ZYallers/zgin)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/ZYallers/zgin.svg?branch=master)](https://travis-ci.org/ZYallers/zgin) 
[![Foundation](https://img.shields.io/badge/Golang-Foundation-green.svg)](http://golangfoundation.org) 
[![GoDoc](https://pkg.go.dev/badge/github.com/ZYallers/zgin?status.svg)](https://pkg.go.dev/github.com/ZYallers/zgin?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/ZYallers/zgin/-/badge.svg)](https://sourcegraph.com/github.com/ZYallers/zgin?badge)
[![Release](https://img.shields.io/github/release/ZYallers/zgin.svg?style=flat-square)](https://github.com/ZYallers/zgin/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/ZYallers/zgin)](https://www.tickgit.com/browse?repo=github.com/ZYallers/zgin)
[![goproxy.cn](https://goproxy.cn/stats/github.com/ZYallers/zgin/badges/download-count.svg)](https://goproxy.cn)

Zgin is a API framework written in Go (Golang). 
An MVCS, Restful, and version control framework based on the [Gin](https://github.com/gin-gonic/gin) framework.
If you need performance and good productivity, you will love zgin.

## Installation
To install zgin package, you need to install Go and set your Go workspace first.

1. The first need Go installed (version 1.11+ is required), then you can use the below Go command to install zgin.
```bash
$ go get -u github.com/ZYallers/zgin
```

2. Import it in your code:
```go 
import "github.com/ZYallers/zgin" 
```

## Quick start

```go
package main

import (
	"fmt"
	"github.com/ZYallers/zgin"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/example/route"
	"github.com/ZYallers/zgin/handler"
	"github.com/ZYallers/zgin/helper/config"
	"github.com/ZYallers/zgin/middleware"
	"github.com/ZYallers/zgin/option"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
	if err := config.ReadFile(); err != nil {
		panic(fmt.Errorf("read config file error: %s", err))
	}
	app := zgin.New(
		option.WithMode(consts.DevMode),
	)
	app.Run(
		handler.WithNoRoute(),
		handler.WithHealth(),
		middleware.WithZapRecovery(),
		middleware.WithZapLogger(),
		handler.WithExpVar(),
		handler.WithPrometheus(),
		handler.WithSwagger(),
		handler.WithPProf(),
		middleware.WithRestCheck(route.Restful),
	)
}

```
run main.go and visit http://0.0.0.0:9010/health (for windows "http://localhost:8080/health") on browser
```bash
$ go run main.go
```

## Zgin v1.stable
- MVCS four-tier architecture support
- Restful interface style support
- API version control and permission custom configuration
- PProf middleware support
- Prometheus middleware support
- Swagger api docs middleware support
- Graceful server shutdown and reload

## Build with jsoniter
Zgin uses encoding/json as default json package but you can change to [jsoniter](https://github.com/json-iterator/go) by build from other tags.
```bash
$ go build -tags=jsoniter .
```

## License
Released under the [MIT License](https://github.com/ZYallers/zgin/blob/master/LICENSE)



